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

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type HandlePaypalWebhookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandlePaypalWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandlePaypalWebhookLogic {
	return &HandlePaypalWebhookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HandlePaypalWebhookLogic) HandlePaypalWebhook(in *dropshipbe.PayPalWebhookRequest) (*dropshipbe.WebhookResponse, error) {

	err := utils.VerifyPayPalSignature(
		l.svcCtx.Config.PayPal.PaypalBaseURL,
		in.Headers,
		in.RawBody,
		l.svcCtx.Config.PayPal.ClientID,
		l.svcCtx.Config.PayPal.Secret,
		l.svcCtx.Config.PayPal.WebhookID,
	)
	if err != nil {
		l.Logger.Errorf("🚨 Invalid paypal signature: %v", err)
		return nil, fmt.Errorf("invalid paypal signature: %v", err)
	}

	var webhookData utils.PayPalWebhookEvent
	if err := json.Unmarshal(in.RawBody, &webhookData); err != nil {
		l.Logger.Errorf("❌ Error Parse JSON Webhook: %v", err)
		return nil, err
	}

	eventID := webhookData.Id
	eventType := webhookData.EventType
	paypalTransactionID := webhookData.Resource.Id
	orderID := webhookData.Resource.CustomId

	if orderID == "" {
		l.Logger.Infof("⚠️ Webhook %s has no custom_id (OrderID). Skip.", eventID)
		return &dropshipbe.WebhookResponse{Success: true}, nil
	}

	l.Logger.Infof("Received Webhook: [%s] for Order: [%s]", eventType, orderID)

	idempotencyKey := constant.PaypalWebhookProcessedKey(eventID)

	isFirstTime, err := l.svcCtx.Redis.SetnxCtx(l.ctx, idempotencyKey, "1")
	if err != nil {
		l.Logger.Errorf("❌ Error connecting to Redis (Setnx): %v", err)
		return nil, err
	}

	if !isFirstTime {
		l.Logger.Infof("🔄 Webhook event %s  already processed. Skip.", eventID)
		return &dropshipbe.WebhookResponse{Success: true, Message: "Already processed"}, nil
	}

	l.svcCtx.Redis.ExpireCtx(l.ctx, idempotencyKey, 7*24*3600) // ttl 7 days

	var processErr error
	switch eventType {
	case "PAYMENT.CAPTURE.COMPLETED", "CHECKOUT.ORDER.APPROVED":
		processErr = l.handlePaymentSuccess(orderID, paypalTransactionID, in.RawBody)
	default:
		l.Logger.Infof("Ignoring event type: %s", eventType)
	}

	if processErr != nil {
		l.svcCtx.Redis.DelCtx(l.ctx, idempotencyKey)
		l.Logger.Errorf("❌ Error  %s: %v", orderID, processErr)
		return nil, processErr
	}

	return &dropshipbe.WebhookResponse{
		Success: true,
		Message: "Webhook processed successfully",
	}, nil
}

func (l *HandlePaypalWebhookLogic) handlePaymentSuccess(orderID string, paypalTransactionID string, rawBody []byte) error {

	// lock with capture api (avoid Race Condition)
	lockKey := constant.OrderPaidLockKey(orderID)

	isFirstTime, err := l.svcCtx.Redis.SetnxCtx(l.ctx, lockKey, "1")
	if err != nil {
		return fmt.Errorf("error setting redis lock: %w", err)
	}

	if !isFirstTime {
		l.Logger.Infof("✅ Order %s already processed. Webhook only records.", orderID)
		return nil
	}
	l.svcCtx.Redis.ExpireCtx(l.ctx, lockKey, 7*24*3600)

	transaction, err := l.svcCtx.EcommerceRepo.GetTransactionByReference(l.ctx, orderID)
	if err != nil {
		l.svcCtx.Redis.DelCtx(l.ctx, lockKey)
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid order id: transaction not found")
		}
		return fmt.Errorf("system error retrieving transaction: %v", err)
	}

	if transaction.Order == nil {
		l.svcCtx.Redis.DelCtx(l.ctx, lockKey)
		return fmt.Errorf("data corruption: order details missing")
	}

	// Double-check DB Level
	if transaction.Status == "completed" || transaction.Order.FinancialStatus == "paid" {
		l.Logger.Infof("Order %s has already been updated to PAID", orderID)
		return nil
	}

	var fullEvent map[string]interface{}
	json.Unmarshal(rawBody, &fullEvent)

	resource, _ := fullEvent["resource"].(map[string]interface{})
	paypalEmail := ""

	if payer, ok := resource["payer"].(map[string]interface{}); ok {
		if email, ok := payer["email_address"].(string); ok {
			paypalEmail = email
		}
	}

	shippingInfo := map[string]string{
		"email":  paypalEmail,
		"source": "paypal_webhook",
	}
	shippingJSON, _ := json.Marshal(shippingInfo)

	l.Logger.Infof("💰 Webhook:Updating order status: %s", transaction.Order.OrderNumber)
	dbErr := l.svcCtx.EcommerceRepo.UpdateTransactionAndOrderStatus(
		l.ctx,
		transaction,
		"completed",
		"paid",
		datatypes.JSON(rawBody),
		map[string]interface{}{
			"customer_email":   paypalEmail,
			"shipping_address": datatypes.JSON(shippingJSON),
		},
	)

	if dbErr != nil {
		l.Logger.Errorf("CRITICAL DB ERROR: Payment collected but DB update failed for order %s: %v", transaction.Order.OrderNumber, dbErr)
		l.svcCtx.Redis.DelCtx(l.ctx, lockKey)
		return dbErr
	}

	// BƯỚC E: BẮN SỰ KIỆN KAFKA ĐỂ GỬI EMAIL
	eventPayload := map[string]interface{}{
		"event_id":       fmt.Sprintf("evt_wh_%s_%d", transaction.Order.OrderNumber, time.Now().UnixNano()),
		"event_type":     "ORDER_PAID",
		"order_id":       transaction.Order.ID,
		"order_number":   transaction.Order.OrderNumber,
		"customer_email": paypalEmail,
		"total_amount":   transaction.Order.TotalPrice,
		"currency":       transaction.Order.Currency,
		"timestamp":      time.Now().Unix(),
		"source":         "webhook",
	}

	msgBytes, _ := json.Marshal(eventPayload)
	if l.svcCtx.KqNotificationPusherClient != nil {
		if err := l.svcCtx.KqNotificationPusherClient.Push(l.ctx, string(msgBytes)); err != nil {
			l.Logger.Errorf("Failed to push Kafka event for order %s: %v", transaction.Order.OrderNumber, err)
		} else {
			l.Logger.Infof("🚀 Webhook successfully published Kafka Event [ORDER_PAID] for order: %s", transaction.Order.OrderNumber)
		}
	} else {
		l.Logger.Errorf("Kafka Client is not initialized. Skipped sending event for order: %s", transaction.Order.OrderNumber)
	}

	return nil
}
