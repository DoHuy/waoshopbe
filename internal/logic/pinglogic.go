package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *dropshipbe.Request) (*dropshipbe.Response, error) {
	l.Logger.Infof("Nhận được yêu cầu Ping từ client")
	return &dropshipbe.Response{
		Pong: "Pong! Hệ thống Dropship backend đang hoạt động ổn định.",
	}, nil
}
