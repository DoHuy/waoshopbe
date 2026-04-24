package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"dropshipbe/common/constant"
	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"
	model "dropshipbe/model/schema"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *dropshipbe.CreateOrderRequest) (*dropshipbe.CreateOrderResponse, error) {

	var variantIDs []uint64
	for _, item := range in.Items {
		variantIDs = append(variantIDs, item.VariantId)
	}

	var variants, err = l.svcCtx.EcommerceRepo.GetVariantsByIDs(l.ctx, variantIDs)
	if len(variants) == 0 {
		return nil, fmt.Errorf("error: no valid product variants found for the provided IDs")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching variants: %v", err)
	}

	variantPriceMap := make(map[uint64]model.Variant)

	for _, v := range variants {
		variantPriceMap[v.ID] = v
	}

	var totalAmount float64
	var orderItems []model.OrderItem

	for _, item := range in.Items {
		v, exists := variantPriceMap[item.VariantId]
		if !exists || v.IsActive == nil || !*v.IsActive {
			return nil, fmt.Errorf("product variant %d is not available", item.VariantId)
		}

		price := v.Price
		quantity := int(item.Quantity)
		lineTotal := price * float64(quantity)
		totalAmount += lineTotal

		orderItems = append(orderItems, model.OrderItem{
			VariantID:   v.ID,
			ProductID:   v.ProductID,
			ProductName: v.Product.Name,
			VariantName: v.Sku,
			Sku:         v.Sku,
			Quantity:    quantity,
			Price:       price,
			Total:       lineTotal,
		})
	}
	l.Logger.Infof("Order Total: %.2f GBP", totalAmount)

	var successfullyDeducted []*dropshipbe.CartItem

	for _, item := range in.Items {
		res, err := utils.DeductInventory(l.ctx, l.svcCtx.Redis, item.VariantId, item.Quantity)
		if err != nil || res != 1 {
			l.Logger.Errorf("Product %s is out of stock or inventory error", item.VariantId)
			utils.RollbackInventory(l.ctx, l.svcCtx.Redis, item.VariantId, item.Quantity)
			return nil, fmt.Errorf("product %d is currently unavailable", item.VariantId)
		}
		successfullyDeducted = append(successfullyDeducted, item)
	}

	email, phone, name, addressStr := "", "", "", ""
	if in.CustomerEmail != nil {
		email = *in.CustomerEmail
	}
	if in.CustomerPhone != nil {
		phone = *in.CustomerPhone
	}
	if in.CustomerName != nil {
		name = *in.CustomerName
	}
	if in.ShippingAddress != nil {
		addressStr = *in.ShippingAddress
	}

	addressMap := map[string]string{
		"recipient_name": name,
		"full_address":   addressStr,
		"method":         in.ShippingMethod,
	}
	addressBytes, _ := json.Marshal(addressMap)

	orderCode := constant.CreateOrderNumber()

	newOrder := model.Order{
		OrderNumber:       orderCode,
		CustomerEmail:     email,
		CustomerPhone:     phone,
		ShippingAddress:   addressBytes,
		TotalPrice:        totalAmount,
		SubtotalPrice:     totalAmount,
		Currency:          "GBP",
		FinancialStatus:   "pending",
		FulfillmentStatus: "unfulfilled",
	}

	err = l.svcCtx.EcommerceRepo.CreateOrder(l.ctx, &newOrder, orderItems)

	if err != nil {
		// If DB error occurs, we must rollback inventory for all successfully deducted items to prevent stock inconsistency
		l.Logger.Errorf("Error Transaction DB: %v", err)
		l.rollbackAllInventory(successfullyDeducted)
		return nil, fmt.Errorf("error system when creating order, please try again later")
	}

	l.Logger.Infof("Calling PayPal for Order: %s, Amount: %.2f", orderCode, totalAmount)

	paypalOrderID, _, err := utils.CreatePayPalOrder(
		l.svcCtx.Config.PayPal.PaypalBaseURL,
		l.svcCtx.Config.PayPal.ClientID,
		l.svcCtx.Config.PayPal.Secret,
		l.svcCtx.Config.PayPal.Mode,
		totalAmount,
	)

	if err != nil {
		// If payment gateway error occurs, we must:
		// 1. Mark order as Canceled in DB (since payment failed)
		// 2. Rollback inventory for all successfully deducted items to prevent stock inconsistency
		l.Logger.Errorf("Error calling PayPal: %v", err)
		l.svcCtx.EcommerceRepo.UpdateOrderStatus(l.ctx, newOrder.ID, "canceled")
		l.rollbackAllInventory(successfullyDeducted)
		return nil, fmt.Errorf("payment gateway is currently down, please try again later")
	}

	newTransaction := model.Transaction{
		OrderID:              newOrder.ID,
		Gateway:              "paypal",
		PaymentMethod:        "paypal_checkout",
		TransactionReference: paypalOrderID,
		Amount:               totalAmount,
		Currency:             "GBP",
		Status:               "pending",
	}

	if err := l.svcCtx.EcommerceRepo.CreateTransaction(l.ctx, &newTransaction); err != nil {
		l.Logger.Errorf("Save Transaction failed for Order %s: %v", orderCode, err)

		// Must rollback inventory for all successfully deducted items to prevent stock inconsistency
		l.svcCtx.EcommerceRepo.UpdateOrderStatus(l.ctx, newOrder.ID, "canceled")

		// rollback all inventory for the items in the order to prevent stock inconsistency
		l.rollbackAllInventory(successfullyDeducted)

		return nil, fmt.Errorf("Error system when processing payment, please contact support")
	}

	return &dropshipbe.CreateOrderResponse{
		LocalOrderId:  orderCode,
		PaypalOrderId: paypalOrderID,
		TotalAmount:   float32(totalAmount),
	}, nil

}

func (l *CreateOrderLogic) rollbackAllInventory(items []*dropshipbe.CartItem) {
	for _, item := range items {
		err := utils.RollbackInventory(context.Background(), l.svcCtx.Redis, item.VariantId, item.Quantity)
		if err != nil {
			l.Logger.Errorf("CRITICAL: Error rolling back inventory for VariantID %s: %v", item.VariantId, err)
		}
	}
}
