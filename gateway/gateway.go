package main

import (
	"flag"
	"log"
	"net/http"

	"dropshipbe/common/middleware"
	"dropshipbe/dropshipbeclient"
	"dropshipbe/gateway/config"
	"dropshipbe/gateway/router"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

var gatewayConfigFile = flag.String("f", "etc/gateway.yaml", "gateway configuration file")

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Please ensure all environment variables in .env are exported: %v", err)
	}

	var c config.GwConfig
	if err := conf.Load(*gatewayConfigFile, &c, conf.UseEnv()); err != nil {
		log.Fatalf("❌ ERROR: Failed to load Gateway YAML file: %v", err)
	}

	logx.MustSetup(c.Log)

	gw := gateway.MustNewServer(c.GatewayConf)
	defer gw.Stop()

	rds := c.Redis.NewRedis()
	rpcConfig := getRpcConfig(c)
	rpcClient := zrpc.MustNewClient(rpcConfig)
	dropshipSvc := dropshipbeclient.NewDropshipbe(rpcClient)

	gw.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, *")
			next(w, r)
		}
	})

	router.RegisterRoutes(gw.Server, dropshipSvc, rds, c)

	middleware.RegisterGatewayMiddlewares(gw.Server, rds)

	setupCorsForGateway(gw.Server, c)

	logx.Infof("🚀 Starting Gateway Server at %s:%d...\n", c.Host, c.Port)
	gw.Start()
}

func getRpcConfig(c config.GwConfig) zrpc.RpcClientConf {
	for _, upstream := range c.Upstreams {
		if upstream.Grpc != nil {
			return *upstream.Grpc
		}
	}
	panic("ERROR: not found any gRPC upstream configuration in gateway.yaml")
}

func setupCorsForGateway(gw *rest.Server, c config.GwConfig) {
	pathMap := make(map[string]bool)

	pathMap["/ws/chat"] = true
	pathMap["/api/v1/webhooks/paypal"] = true

	for _, upstream := range c.Upstreams {
		for _, mapping := range upstream.Mappings {
			pathMap[mapping.Path] = true
		}
	}

	var optionsRoutes []rest.Route
	for path := range pathMap {
		optionsRoutes = append(optionsRoutes, rest.Route{
			Method: http.MethodOptions,
			Path:   path,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
		})
	}

	gw.AddRoutes(optionsRoutes)
}
