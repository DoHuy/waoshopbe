package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductFaqsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductFaqsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductFaqsLogic {
	return &GetProductFaqsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- FAQs ---
func (l *GetProductFaqsLogic) GetProductFaqs(in *dropshipbe.GetProductFaqsRequest) (*dropshipbe.FaqListResponse, error) {
	faqs, err := l.svcCtx.EcommerceRepo.GetProductFaqs(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var faqItems []*dropshipbe.Faq
	for _, f := range faqs {
		faqItems = append(faqItems, &dropshipbe.Faq{
			Id:       f.ID,
			Question: f.Question,
			Answer:   f.Answer,
		})
	}

	return &dropshipbe.FaqListResponse{Faqs: faqItems}, err
}
