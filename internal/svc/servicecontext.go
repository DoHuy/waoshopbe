package svc

import (
	"context"
	"dropshipbe/internal/config"
	"dropshipbe/model/repository"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config                       config.Config
	S3Client                     *s3.Client
	PresignClient                *s3.PresignClient
	DB                           *gorm.DB
	Redis                        *redis.Redis
	EcommerceRepo                repository.EcommerceRepository
	KqOrderPusherClient          *kq.Pusher
	KqNotificationPusherClient   *kq.Pusher
	KqEmailMarketingPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.R2.AccountID)

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.R2.AccessKey,
			c.R2.SecretKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Cannot load R2 configuration: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	presignClient := s3.NewPresignClient(s3Client)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          //  color
		},
	)

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Ho_Chi_Minh", c.DB.Host, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.Port, c.DB.SSLMode)), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Cannot connect to Database: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(c.DB.MaxOpenConns)
		sqlDB.SetMaxIdleConns(c.DB.MaxIdleConns)
	}

	dropShipCache := cache.New(
		c.CacheConf,
		syncx.NewSingleFlight(),
		cache.NewStat("dropship_cache"),
		gorm.ErrRecordNotFound,
		cache.WithExpiry(time.Duration(c.CacheTTL)*time.Minute),
	)

	flushDuration := time.Duration(c.KqPusherConf.FlushInterval) * time.Millisecond

	pushOptions := []kq.PushOption{
		kq.WithChunkSize(c.KqPusherConf.ChunkSize),
		kq.WithFlushInterval(flushDuration),
	}

	var rds *redis.Redis
	if len(c.CacheConf) > 0 {
		rds = redis.New(c.CacheConf[0].Host, redis.WithPass(c.CacheConf[0].Pass))
	} else {
		log.Fatalf("Cannot find Cache/Redis configuration")
	}

	return &ServiceContext{
		Config:                       c,
		DB:                           db,
		Redis:                        rds,
		S3Client:                     s3Client,
		PresignClient:                presignClient,
		EcommerceRepo:                repository.NewEcommerceRepository(db, dropShipCache),
		KqOrderPusherClient:          kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.OrderTopic, pushOptions...),
		KqNotificationPusherClient:   kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.NotificationTopic, pushOptions...),
		KqEmailMarketingPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.EmailMarketingTopic, pushOptions...),
	}
}
