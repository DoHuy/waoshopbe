package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNewBlogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNewBlogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNewBlogLogic {
	return &CreateNewBlogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateNewBlogLogic) CreateNewBlog(in *dropshipbe.CreateNewBlogRequest) (*dropshipbe.Blog, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.Blog{}, nil
}
