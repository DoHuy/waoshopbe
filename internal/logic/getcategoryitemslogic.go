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

type GetCategoryItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCategoryItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryItemsLogic {
	return &GetCategoryItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCategoryItemsLogic) GetCategoryItems(in *dropshipbe.DefaultRequest) (*dropshipbe.CategoryListResponse, error) {
	// todo: add your logic here and delete this line
	categories, err := l.svcCtx.EcommerceRepo.GetCategoryItems(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var categoryItems []*dropshipbe.Category

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, category := range categories {

		if category.ImageURL != "" {
			presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(category.ImageURL),
			}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
			if err != nil {
				logx.Errorf("Lỗi khi tạo presigned URL cho ảnh category: %v", err)
				return nil, err
			}
			category.ImageURL = presignedURL.URL
		}
		categoryItems = append(categoryItems, &dropshipbe.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CountryCode: category.CountryCode,
			Slug:        category.Slug,
			ImageUrl:    category.ImageURL,
			ParentId:    category.ParentID,
			IsActive:    *category.IsActive,
		})
	}

	return &dropshipbe.CategoryListResponse{
		Categories: categoryItems,
	}, nil
}
