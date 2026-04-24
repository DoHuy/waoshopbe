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

type GetSliderItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSliderItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSliderItemsLogic {
	return &GetSliderItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- UI Items (Sliders, Categories, Banners) ---
func (l *GetSliderItemsLogic) GetSliderItems(in *dropshipbe.DefaultRequest) (*dropshipbe.SliderListResponse, error) {
	// todo: add your logic here and delete this line
	sliders, err := l.svcCtx.EcommerceRepo.GetSliderItems(l.ctx, in)
	if err != nil {
		logx.Errorf("Lỗi khi lấy slider items: %v", err)
		return nil, err
	}

	var sliderItems []*dropshipbe.Slider
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, s := range sliders {
		// image to presigned URL
		if s.ImageURL != "" {
			presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(s.ImageURL),
			}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
			if err != nil {
				logx.Errorf("Lỗi khi tạo presigned URL cho slider image: %v", err)
				s.ImageURL = ""
			} else {
				s.ImageURL = presignedURL.URL
			}
		}

		sliderItems = append(sliderItems, &dropshipbe.Slider{
			ImgSrc:      s.ImageURL,
			Title:       s.Title,
			SubText:     s.SubText,
			Description: s.Description,
		})
	}

	return &dropshipbe.SliderListResponse{Sliders: sliderItems}, nil

}
