package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dropshipbe/common/constant"
	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"
	model "dropshipbe/model/schema"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CaptureOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptureOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptureOrderLogic {
	return &CaptureOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CaptureOrderLogic) CaptureOrder(in *dropshipbe.CaptureOrderRequest) (*dropshipbe.CaptureOrderResponse, error) {

	transaction, err := l.svcCtx.EcommerceRepo.GetTransactionByReference(l.ctx, in.PaypalOrderId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid paypal order id: transaction not found")
		}
		return nil, fmt.Errorf("system error retrieving transaction: %v", err)
	}

	if transaction.Order == nil {
		return nil, fmt.Errorf("data corruption: order details missing")
	}

	if transaction.Status == "completed" || transaction.Order.FinancialStatus == "paid" {
		return &dropshipbe.CaptureOrderResponse{
			Success: true,
			Status:  "COMPLETED",
			Message: "This order has already been successfully paid.",
		}, nil
	}

	l.Logger.Infof("Processing capture for PayPal Order ID: %s", in.PaypalOrderId)

	status, rawResponse, err := utils.CapturePayPalOrder(
		l.svcCtx.Config.PayPal.PaypalBaseURL,
		l.svcCtx.Config.PayPal.ClientID,
		l.svcCtx.Config.PayPal.Secret,
		in.PaypalOrderId,
	)

	rawJSON, _ := json.Marshal(rawResponse)

	if err != nil || status != "COMPLETED" {
		l.Logger.Errorf("Capture failed. Status: %s, Error: %v", status, err)

		l.svcCtx.EcommerceRepo.UpdateTransactionAndOrderStatus(l.ctx, transaction, "failed", "canceled", datatypes.JSON(rawJSON), nil)
		l.rollbackInventory(transaction.Order.OrderItems)

		return &dropshipbe.CaptureOrderResponse{
			Success: false,
			Status:  status,
			Message: "Transaction declined. Please check your payment method or PayPal account.",
		}, nil
	}

	l.Logger.Infof("💰 Successfully captured payment from PayPal for order: %s", transaction.Order.OrderNumber)

	lockKey := constant.OrderPaidLockKey(transaction.Order.OrderNumber)

	isFirstTime, redisErr := l.svcCtx.Redis.SetnxCtx(l.ctx, lockKey, "1")
	if redisErr != nil {
		l.Logger.Errorf("Redis Lock error: %v", redisErr)
		isFirstTime = true
	}

	if !isFirstTime {
		l.Logger.Infof("✅ Đơn hàng %v đã được Webhook xử lý lưu DB & Kafka. API Capture chỉ trả kết quả.", transaction.Order.OrderNumber)
		return &dropshipbe.CaptureOrderResponse{
			Success: true,
			Status:  "COMPLETED",
			Message: "Payment successful. Your order is being prepared.",
		}, nil
	}

	l.svcCtx.Redis.ExpireCtx(l.ctx, lockKey, 7*24*3600)

	paypalEmail := ""
	if payer, ok := rawResponse["payer"].(map[string]interface{}); ok {
		if email, ok := payer["email_address"].(string); ok {
			paypalEmail = email
		}
	}

	var shippingName, shippingAddressStr string
	if purchaseUnits, ok := rawResponse["purchase_units"].([]interface{}); ok && len(purchaseUnits) > 0 {
		if firstUnit, ok := purchaseUnits[0].(map[string]interface{}); ok {
			if shipping, ok := firstUnit["shipping"].(map[string]interface{}); ok {
				if nameObj, ok := shipping["name"].(map[string]interface{}); ok {
					shippingName, _ = nameObj["full_name"].(string)
				}
				if addr, ok := shipping["address"].(map[string]interface{}); ok {
					shippingAddressStr = fmt.Sprintf("%v, %v, %v, %v, %v",
						addr["address_line_1"], addr["admin_area_2"], addr["admin_area_1"], addr["postal_code"], addr["country_code"],
					)
				}
			}
		}
	}

	shippingInfo := map[string]string{
		"recipient_name": shippingName,
		"full_address":   shippingAddressStr,
		"email":          paypalEmail,
		"source":         "paypal_capture",
	}
	shippingJSON, _ := json.Marshal(shippingInfo)

	dbErr := l.svcCtx.EcommerceRepo.UpdateTransactionAndOrderStatus(l.ctx, transaction, "completed", "paid", datatypes.JSON(rawJSON), map[string]interface{}{
		"customer_email":   paypalEmail,
		"shipping_address": datatypes.JSON(shippingJSON),
	})

	if dbErr != nil {
		l.Logger.Errorf("CRITICAL DB ERROR: Payment collected but DB update failed for order %s: %v", transaction.Order.OrderNumber, dbErr)
		l.svcCtx.Redis.DelCtx(l.ctx, lockKey)

		return &dropshipbe.CaptureOrderResponse{
			Success: false,
			Status:  "ERROR",
			Message: "Payment succeeded but system encountered an error updating your order. Support will contact you.",
		}, dbErr
	}

	eventPayload := map[string]interface{}{
		"event_id":       fmt.Sprintf("evt_%s_%d", transaction.Order.OrderNumber, time.Now().UnixNano()),
		"event_type":     "ORDER_PAID",
		"order_id":       transaction.Order.ID,
		"order_number":   transaction.Order.OrderNumber,
		"customer_email": paypalEmail,
		"customer_name":  shippingName,
		"total_amount":   transaction.Order.TotalPrice,
		"currency":       transaction.Order.Currency,
		"timestamp":      time.Now().Unix(),
		"source":         "api_capture",
	}

	msgBytes, err := json.Marshal(eventPayload)
	if err != nil {
		l.Logger.Errorf("Failed to Marshal Kafka Event for order %s: %v", transaction.Order.OrderNumber, err)
	} else {
		if l.svcCtx.KqNotificationPusherClient != nil {
			err = l.svcCtx.KqNotificationPusherClient.Push(l.ctx, string(msgBytes))
			if err != nil {
				l.Logger.Errorf("Failed to push Kafka event for order %s: %v", transaction.Order.OrderNumber, err)
			} else {
				l.Logger.Infof("🚀 Successfully published Kafka Event [ORDER_PAID] for order: %s", transaction.Order.OrderNumber)
			}
		} else {
			l.Logger.Errorf("Kafka Client is not initialized in ServiceContext. Skipped sending event for order: %s", transaction.Order.OrderNumber)
		}
	}

	return &dropshipbe.CaptureOrderResponse{
		Success: true,
		Status:  "COMPLETED",
		Message: "Payment successful. Your order is being prepared.",
	}, nil
}

func (l *CaptureOrderLogic) rollbackInventory(items []model.OrderItem) {
	for _, item := range items {
		err := utils.RollbackInventory(l.ctx, l.svcCtx.Redis, item.VariantID, int32(item.Quantity))
		if err != nil {
			l.Logger.Errorf("CRITICAL: Failed to rollback inventory for VariantID %d: %v", item.VariantID, err)
		}
	}
}
