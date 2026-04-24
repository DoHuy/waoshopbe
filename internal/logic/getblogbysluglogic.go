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

type GetBlogBySlugLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlogBySlugLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlogBySlugLogic {
	return &GetBlogBySlugLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBlogBySlugLogic) GetBlogBySlug(in *dropshipbe.GetBlogBySlugRequest) (*dropshipbe.BlogDetailResponse, error) {
	// todo: add your logic here and delete this line
	blogPost, err := l.svcCtx.EcommerceRepo.GetBlogBySlug(l.ctx, in)
	if err != nil {
		logx.Errorf("Lỗi khi lấy blog theo slug: %v", err)
		return nil, err
	}

	blogPostResponse := new(dropshipbe.Blog)

	if blogPost.ImageURL != "" {
		presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(l.ctx, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(blogPost.ImageURL),
		}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
		if err != nil {
			logx.Errorf("Lỗi khi tạo presigned URL cho ảnh blog: %v", err)
			return nil, err
		}
		blogPost.ImageURL = presignedURL.URL
	}

	blogPostResponse.Id = blogPost.ID
	blogPostResponse.Title = blogPost.Title
	blogPostResponse.Content = blogPost.Content
	blogPostResponse.Slug = blogPost.Slug

	blogPostResponse.Alt = blogPost.ImageAlt
	blogPostResponse.ImageUrl = blogPost.ImageURL
	blogPostResponse.Category = &dropshipbe.BlogCategory{
		Title: blogPost.Category.Name,
		Slug:  blogPost.Category.Slug,
	}

	return &dropshipbe.BlogDetailResponse{
		Blog: blogPostResponse,
	}, nil
}
