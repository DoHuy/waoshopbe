package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"dropshipbe/mq/internal/config"
	"dropshipbe/mq/internal/consumer"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/mq.yaml", "the config file")

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var c config.Config
	err := conf.Load(*configFile, &c, conf.UseEnv())
	if err != nil {
		log.Fatalf("❌ Error: Load file YAML failed: %v", err)
	}

	logx.MustSetup(c.Log)

	sg := service.NewServiceGroup()
	defer sg.Stop()

	notificationLogic := consumer.NewNotificationConsumer(ctx, c)
	notificationQueue := kq.MustNewQueue(c.NotificationConsumer, notificationLogic)
	sg.Add(notificationQueue)

	logx.Info("🚀 MQ Service is running...")
	sg.Start()
}
