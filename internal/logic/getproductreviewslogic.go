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

type GetProductReviewsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductReviewsLogic {
	return &GetProductReviewsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Reviews ---
func (l *GetProductReviewsLogic) GetProductReviews(in *dropshipbe.GetProductReviewsRequest) (*dropshipbe.ReviewSummary, error) {
	// todo: add your logic here and delete this line
	productReviews, err := l.svcCtx.EcommerceRepo.GetProductReviews(l.ctx, in)
	if err != nil {
		return nil, err
	}

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var reviews []*dropshipbe.ReviewItem
	for _, r := range productReviews {

		// images to presigned URLs
		// Assuming r.Media.Images is a slice of image keys in R2

		var mediaURLs []string
		var videoURLs []string
		var authorAvatar string

		if r.AuthorAvatar != "" {
			presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(r.AuthorAvatar),
			}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
			if err != nil {
				logx.Errorf("Lỗi khi tạo presigned URL cho review author avatar image: %v", err)
			}

			authorAvatar = presignedURL.URL
		}

		if r.Media != nil {
			for _, imageKey := range r.Media.Images {
				presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
					Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
					Key:    aws.String(imageKey),
				}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
				if err != nil {
					logx.Errorf("Lỗi khi tạo presigned URL cho review image: %v", err)
					continue // Skip this image if there's an error
				}
				mediaURLs = append(mediaURLs, presignedURL.URL)
			}
			r.Media.Images = mediaURLs // Update the Media field with presigned URLs

			for _, videoKey := range r.Media.Videos {
				presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
					Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
					Key:    aws.String(videoKey),
				}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
				if err != nil {
					logx.Errorf("Lỗi khi tạo presigned URL cho review video: %v", err)
					continue // Skip this video if there's an error
				}
				videoURLs = append(videoURLs, presignedURL.URL)
			}
			r.Media.Videos = videoURLs // Update the Media field with presigned URLs

			reviews = append(reviews, &dropshipbe.ReviewItem{
				Id:       r.ID,
				Verified: r.IsVerified,
				Rating:   int32(r.Rating),
				Comment:  r.Content,
				Images:   r.Media.Images,
				Videos:   r.Media.Videos,
				Name:     r.AuthorName,
				Date:     r.CreatedAt.Format(time.RFC3339),
				Avatar:   authorAvatar,
			})
		}
	}

	var mapRating = map[string]int32{"1": 0, "2": 0, "3": 0, "4": 0, "5": 0}
	if len(reviews) > 0 {
		var totalRating int
		for _, review := range reviews {
			totalRating += int(review.Rating)
			switch review.Rating {
			case 1:
				mapRating["1"] += 1
			case 2:
				mapRating["2"] += 1
			case 3:
				mapRating["3"] += 1
			case 4:
				mapRating["4"] += 1
			case 5:
				mapRating["5"] += 1
			}
		}
		ratingAverage := float64(totalRating) / float64(len(reviews))

		return &dropshipbe.ReviewSummary{Reviews: reviews, RatingAverage: float32(ratingAverage), RatingCount: int32(len(reviews)), Rating: mapRating}, nil
	}

	return &dropshipbe.ReviewSummary{Reviews: reviews, RatingAverage: 0, RatingCount: 0, Rating: nil}, nil
}
