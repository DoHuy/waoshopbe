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

type GetVideoBannerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoBannerLogic {
	return &GetVideoBannerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoBannerLogic) GetVideoBanner(in *dropshipbe.DefaultRequest) (*dropshipbe.Banner, error) {
	// todo: add your logic here and delete this line
	socialVideo, err := l.svcCtx.EcommerceRepo.GetVideoBanner(l.ctx, in)
	if err != nil {
		logx.Errorf("Lỗi khi lấy video banner: %v", err)
		return nil, err
	}

	// Chuyển đổi dữ liệu từ model sang response
	if socialVideo == nil {
		return &dropshipbe.Banner{}, nil
	}

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// presigned video
	if socialVideo.VideoURL != "" {
		presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(socialVideo.VideoURL),
		}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
		if err != nil {
			logx.Errorf("Lỗi khi tạo presigned URL cho video banner: %v", err)
			socialVideo.VideoURL = "" // fallback nếu có lỗi
		} else {
			socialVideo.VideoURL = presignedURL.URL
		}
	}

	return &dropshipbe.Banner{
		Id:          socialVideo.ID,
		VideoUrl:    &socialVideo.VideoURL,
		Alt:         socialVideo.Alt,
		Title:       socialVideo.Title,
		Description: socialVideo.Description,
		VideoType:   &socialVideo.VideoType,
	}, nil
}
