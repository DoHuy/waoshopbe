package logic

import (
	"context"
	"fmt"
	"time"

	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGuestTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGuestTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGuestTokenLogic {
	return &GetGuestTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Guest token for chatbot
func (l *GetGuestTokenLogic) GetGuestToken(in *dropshipbe.GetGuestTokenRequest) (*dropshipbe.GetGuestTokenResponse, error) {

	guestSessionID := fmt.Sprintf("guest_%d", time.Now().Unix())

	signedToken, expireAt, err := utils.GenerateGuestToken(guestSessionID, l.svcCtx.Config.Jwt)

	return &dropshipbe.GetGuestTokenResponse{
		Token:    signedToken,
		ExpireAt: expireAt,
	}, err
}
