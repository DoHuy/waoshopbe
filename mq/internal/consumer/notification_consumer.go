package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"dropshipbe/common/kafka_events"
	"dropshipbe/common/notify"
	"dropshipbe/mq/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationConsumer struct {
	Config    config.Config
	ctx       context.Context
	sesClient *notify.SESClient
	logx.Logger
}

func NewNotificationConsumer(ctx context.Context, c config.Config) *NotificationConsumer {

	ses, err := notify.NewSESClient(
		context.Background(),
		c.Email.Region,
		c.Email.AccessKey,
		c.Email.SecretKey,
		c.Email.FromAddress,
	)

	if err != nil {
		logx.Must(err)
	}

	return &NotificationConsumer{
		ctx:       ctx,
		Config:    c,
		sesClient: ses,
		Logger:    logx.WithContext(ctx),
	}
}

func (c *NotificationConsumer) Consume(ctx context.Context, key, val string) error {
	c.Logger.Infof("📥 Receive from Kafka: %s", val)

	var payload kafka_events.OrderEventPayload
	if err := json.Unmarshal([]byte(val), &payload); err != nil {
		c.Logger.Errorf("Error parsing JSON event: %v", err)
		return nil
	}

	switch payload.EventType {
	case "ORDER_PAID":
		return c.handleOrderPaid(ctx, payload)
	default:
		c.Logger.Infof("Without handler for event type: %s", payload.EventType)
		return nil
	}
}

func (c *NotificationConsumer) handleOrderPaid(ctx context.Context, payload kafka_events.OrderEventPayload) error {

	// ==========================================
	// Send to ADMIN
	// ==========================================
	teleMsg := fmt.Sprintf(
		"💸 <b>NEW ORDER PAID</b>\n"+
			"Order ID: %s\nCustomer: %s\nRevenue: %.2f %s",
		payload.OrderNumber, payload.CustomerName, payload.TotalAmount, payload.Currency,
	)

	err := notify.SendTelegramMessage(c.Config.Telegram.BotToken, c.Config.Telegram.ChatID, teleMsg)
	if err != nil {
		c.Logger.Errorf("❌ Error sending Telegram for order %s: %v", payload.OrderNumber, err)
	}

	if payload.CustomerEmail != "" {
		if c.sesClient == nil {
			c.Logger.Errorf("❌ Error: SES client is not initialized for customer email: %s", payload.CustomerEmail)
		} else {
			subject := fmt.Sprintf("Thank you for your order %s!", payload.OrderNumber)
			body := fmt.Sprintf(`
                <div style="font-family: Arial, sans-serif; color: #333;">
                    <h2>Hi %s,</h2>
                    <p>We have successfully received your payment of <b>%.2f %s</b>.</p>
                    <p>Your order <strong>%s</strong> is currently being prepared for shipment.</p>
                    <br>
                    <p>Thank you for shopping with us!</p>
                </div>`,
				payload.CustomerName, payload.TotalAmount, payload.Currency, payload.OrderNumber,
			)

			err := c.sesClient.SendEmail(ctx, payload.CustomerEmail, subject, body)
			if err != nil {
				c.Logger.Errorf("❌ Error sending Email via SES for %s: %v", payload.CustomerEmail, err)
			} else {
				c.Logger.Infof("📧 Successfully sent invoice email to customer: %s", payload.CustomerEmail)
			}
		}
	} else {
		c.Logger.Infof("⚠️ Order %s has no customer email, skipping email sending.", payload.OrderNumber)
	}
	c.Logger.Infof("✅ Successfully processed notification for order: %s", payload.OrderNumber)

	return nil
}
