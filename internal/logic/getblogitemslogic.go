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

type GetBlogItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlogItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlogItemsLogic {
	return &GetBlogItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Blogs ---
func (l *GetBlogItemsLogic) GetBlogItems(in *dropshipbe.DefaultRequest) (*dropshipbe.BlogListResponse, error) {
	// todo: add your logic here and delete this line
	blogPosts, err := l.svcCtx.EcommerceRepo.GetBlogItems(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var blogItemsResponse []*dropshipbe.Blog
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, blogPost := range blogPosts {
		blogPostResponse := new(dropshipbe.Blog)

		if blogPost.ImageURL != "" {
			presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
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
		blogItemsResponse = append(blogItemsResponse, blogPostResponse)
	}

	return &dropshipbe.BlogListResponse{
		Blogs: blogItemsResponse,
	}, nil
}
