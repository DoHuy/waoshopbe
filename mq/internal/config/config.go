package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf

	NotificationConsumer kq.KqConf
	FulfillmentConsumer  kq.KqConf `json:",optional"`
	Telegram             struct {
		BotToken string
		ChatID   int64
	}
	Email struct {
		Region      string
		AccessKey   string
		SecretKey   string
		FromAddress string
	}
}
