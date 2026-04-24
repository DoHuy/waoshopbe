package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/config"
	"dropshipbe/internal/server"
	"dropshipbe/internal/svc"
)

var configFile = flag.String("f", "etc/dropshipbe.yaml", "the config file")

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Error loading .env file: %v", err)
	}

	var c config.Config
	err := conf.Load(*configFile, &c, conf.UseEnv())
	if err != nil {
		log.Fatalf("❌ Error: Failed to load YAML file: %v", err)
	}

	logx.MustSetup(c.Log)
	ctx := svc.NewServiceContext(c)
	logx.Info("✅ Logging initialized successfully.")

	// Pre-warming inventory data to Redis at startup
	logx.Info("🔄 Pre-warming inventory data to Redis...")
	if err := utils.LoadAllInventoryToRedis(context.Background(), ctx.Redis, ctx.DB); err != nil {
		logx.WithContext(context.Background()).Errorf("❌ Error pre-warming inventory data: %v", err)
		panic(fmt.Sprintf("❌ Error: pre-warming inventory data failed: %v", err))
	}
	logx.Info("✅ Pre-warming inventory data completed successfully.")

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dropshipbe.RegisterDropshipbeServer(grpcServer, server.NewDropshipbeServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
