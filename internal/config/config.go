package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	CacheConf cache.CacheConf
	CacheTTL  int64
	R2        struct {
		AccountID      string
		AccessKey      string
		SecretKey      string
		BucketName     string
		LinkExpiration int
	}
	DB struct {
		Host         string
		Port         int
		User         string
		Password     string
		DBName       string
		SSLMode      string
		MaxOpenConns int
		MaxIdleConns int
	}
	KqPusherConf struct {
		Brokers             []string
		OrderTopic          string
		NotificationTopic   string
		EmailMarketingTopic string
		ChunkSize           int `json:",default=500"`
		FlushInterval       int `json:",default=100"` // Flush interval in milliseconds
	}
	PayPal struct {
		Mode          string
		ClientID      string
		Secret        string
		PaypalBaseURL string
		WebhookID     string `json:",optional"`
	}
	OpenAI struct {
		APIKey string
	}
	Jwt struct {
		Secret      string
		ExpireHours int64
	}
}
