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

type GetSocialProductVideosLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSocialProductVideosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSocialProductVideosLogic {
	return &GetSocialProductVideosLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Media/Galleries ---
func (l *GetSocialProductVideosLogic) GetSocialProductVideos(in *dropshipbe.GetSocialProductVideoRequest) (*dropshipbe.GalleryListResponse, error) {
	// todo: add your logic here and delete this line
	socialVideos, err := l.svcCtx.EcommerceRepo.GetSocialProductVideos(l.ctx, in)
	if err != nil {
		return nil, err
	}
	var videoItems []*dropshipbe.Gallery
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	for _, v := range socialVideos {
		// presigned video
		if v.VideoURL != "" {
			// Assuming v.VideoURL is the key in R2
			presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(v.VideoURL),
			}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
			if err != nil {
				logx.Errorf("Lỗi khi tạo presigned URL cho social video: %v", err)
				v.VideoURL = "" // fallback nếu có lỗi
			} else {
				v.VideoURL = presignedURL.URL
			}
		}

		videoItems = append(videoItems, &dropshipbe.Gallery{
			VideoUrl:  v.VideoURL,
			AltText:   v.AltText,
			MediaType: v.MediaType,
			Highlight: v.Highlight,
			Position:  int32(v.Position),
		})
	}

	return &dropshipbe.GalleryListResponse{}, nil
}
