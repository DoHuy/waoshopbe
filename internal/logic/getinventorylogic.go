package logic

import (
	"context"
	"fmt"
	"strconv"

	"dropshipbe/common/constant"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryLogic {
	return &GetInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Inventory ---
func (l *GetInventoryLogic) GetInventory(in *dropshipbe.GetInventoryRequest) (*dropshipbe.GetInventoryResponse, error) {
	// todo: add your logic here and delete this line
	type Variant struct {
		ID uint64
	}
	var variants []Variant

	// Truy vấn DB lấy các variant đang active của sản phẩm này
	// Giả định bạn đang dùng GORM trong l.svcCtx.DB
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Table("variants").
		Select("id").
		Where("product_id = ? AND is_active = ?", in.ProductId, true).
		Find(&variants).Error; err != nil {

		l.Logger.Errorf("Lỗi khi lấy danh sách variant từ DB cho product %d: %v", in.ProductId, err)
		return nil, fmt.Errorf("không thể lấy thông tin biến thể của sản phẩm")
	}

	if len(variants) == 0 {
		return &dropshipbe.GetInventoryResponse{
			ProductId:     in.ProductId,
			TotalStock:    0,
			VariantStocks: nil,
		}, nil
	}

	// 2. Khởi tạo biến lưu kết quả
	var variantStocks []*dropshipbe.VariantInventory
	var totalStock int32 = 0

	// 3. Lấy tồn kho từng biến thể từ Redis
	for _, v := range variants {
		// Tạo key Redis theo format bạn đã định nghĩa trong utils
		key := constant.ProductVariantStock(strconv.FormatUint(v.ID, 10))

		// Đọc từ Redis
		valStr, err := l.svcCtx.Redis.GetCtx(l.ctx, key)
		if err != nil {
			l.Logger.Infof("Không tìm thấy tồn kho cho variant %d trong Redis", v.ID)
			continue
		}

		// Chuyển chuỗi trả về từ Redis sang số nguyên
		stock, err := strconv.Atoi(valStr)
		if err != nil {
			l.Logger.Errorf("Lỗi ép kiểu tồn kho từ Redis cho variant %d: %v", v.ID, err)
			continue
		}

		// Đưa vào danh sách trả về
		variantStocks = append(variantStocks, &dropshipbe.VariantInventory{
			VariantId:     v.ID,
			StockQuantity: int32(stock),
		})

		// Cộng dồn tổng kho
		totalStock += int32(stock)
	}

	// 4. Trả kết quả về cho Client
	return &dropshipbe.GetInventoryResponse{
		ProductId:     in.ProductId,
		TotalStock:    totalStock,
		VariantStocks: variantStocks,
	}, nil
}
