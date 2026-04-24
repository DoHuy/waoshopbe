package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/gateway"
)

type JwtConfig struct {
	Secret      string
	ExpireHours int64
}

type GwConfig struct {
	gateway.GatewayConf
	Redis redis.RedisConf
	Jwt   JwtConfig
}
