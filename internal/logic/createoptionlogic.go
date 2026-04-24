package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOptionLogic {
	return &CreateOptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Options ---
func (l *CreateOptionLogic) CreateOption(in *dropshipbe.CreateOptionRequest) (*dropshipbe.Option, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.Option{}, nil
}
