package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductFaqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductFaqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductFaqLogic {
	return &CreateProductFaqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductFaqLogic) CreateProductFaq(in *dropshipbe.CreateProductFaqRequest) (*dropshipbe.Faq, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.Faq{}, nil
}
