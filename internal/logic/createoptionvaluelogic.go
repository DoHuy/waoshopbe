package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOptionValueLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOptionValueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOptionValueLogic {
	return &CreateOptionValueLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOptionValueLogic) CreateOptionValue(in *dropshipbe.CreateOptionValueRequest) (*dropshipbe.OptionValue, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.OptionValue{}, nil
}
