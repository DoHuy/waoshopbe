package logic

import (
	"context"
	"time"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetBannerItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBannerItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBannerItemsLogic {
	return &GetBannerItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBannerItemsLogic) GetBannerItems(in *dropshipbe.DefaultRequest) (*dropshipbe.BannerListResponse, error) {
	// todo: add your logic here and delete this line
	bannerItems, err := l.svcCtx.EcommerceRepo.GetBannerItems(l.ctx, in)
	if err != nil {
		logx.Errorf("Lỗi khi lấy banner: %v", err)
		return nil, err
	}

	// Chuyển đổi dữ liệu từ model sang response
	expirationDuration := time.Duration(l.svcCtx.Config.R2.LinkExpiration) * time.Minute
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var bannerItemsResponse []*dropshipbe.Banner
	for _, b := range bannerItems {

		// image to presigned URL
		if b.ImageURL != "" {
			presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(b.ImageURL),
			}, s3.WithPresignExpires(expirationDuration))
			if err != nil {
				logx.Errorf("Lỗi khi tạo presigned URL cho banner image: %v", err)
				presignedURL.URL = "" // fallback nếu có lỗi
			}
			b.ImageURL = presignedURL.URL
		}

		if b.VideoURL != "" {
			presignedVideoURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(b.VideoURL),
			}, s3.WithPresignExpires(expirationDuration))
			if err != nil {
				logx.Errorf("Lỗi khi tạo presigned URL cho banner video: %v", err)
				presignedVideoURL.URL = "" // fallback nếu có lỗi
			}
			b.VideoURL = presignedVideoURL.URL
		}

		bannerItemsResponse = append(bannerItemsResponse, &dropshipbe.Banner{
			Id:          b.ID,
			ImageUrl:    b.ImageURL,
			VideoUrl:    &b.VideoURL,
			Alt:         b.Alt,
			Description: b.Description,
			VideoType:   &b.VideoType,
			Title:       b.Title,
		})
	}

	return &dropshipbe.BannerListResponse{
		Banners: bannerItemsResponse,
	}, nil
}
