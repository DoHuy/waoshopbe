package utils

import (
	"context"
	"dropshipbe/common/constant"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

const deductStockScript = `
local stockKey = KEYS[1]
local decrementAmount = tonumber(ARGV[1])

local currentStock = tonumber(redis.call('get', stockKey))
if currentStock == nil then
    return -1
end

if currentStock >= decrementAmount then
    redis.call('decrby', stockKey, decrementAmount)
    return 1
else
    return 0
end
`

var deductStockLua = redis.NewScript(deductStockScript)

func DeductInventory(ctx context.Context, rds *redis.Redis, ProductVariantID uint64, quantity int32) (int, error) {
	key := constant.ProductVariantStock(fmt.Sprintf("%d", ProductVariantID))

	resp, err := rds.ScriptRunCtx(ctx, deductStockLua, []string{key}, quantity)
	if err != nil {
		return -1, fmt.Errorf("Error executing Lua script: %v", err)
	}

	result, ok := resp.(int64)
	if !ok {
		return -1, fmt.Errorf("Error: unable to cast result to int64, got: %v", resp)
	}

	return int(result), nil
}

func RollbackInventory(ctx context.Context, rds *redis.Redis, productVariantID uint64, quantity int32) error {
	key := constant.ProductVariantStock(fmt.Sprintf("%d", productVariantID))

	_, err := rds.IncrbyCtx(ctx, key, int64(quantity))
	if err != nil {
		return fmt.Errorf("Error rolling back inventory for variant %d: %v", productVariantID, err)
	}

	return nil
}

func LoadAllInventoryToRedis(ctx context.Context, rds *redis.Redis, db *gorm.DB) error {
	type Variant struct {
		ID            uint64
		StockQuantity int
	}

	var variants []Variant

	if err := db.WithContext(ctx).Table("variants").Where("is_active = ?", true).Select("id, stock_quantity").Find(&variants).Error; err != nil {
		return fmt.Errorf("Error when querying DB for inventory: %v", err)
	}

	err := rds.PipelinedCtx(ctx, func(pipe redis.Pipeliner) error {
		for _, v := range variants {
			variantID := strconv.FormatUint(v.ID, 10)
			key := constant.ProductVariantStock(variantID)

			pipe.Set(ctx, key, v.StockQuantity, 0)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error when loading inventory via Pipeline: %v", err)
	}

	return nil
}

func GetInventory(ctx context.Context, rds *redis.Redis, productVariantID uint64) (int, error) {

	key := constant.ProductVariantStock(fmt.Sprintf("%d", productVariantID))

	stockStr, err := rds.GetCtx(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("Error fetching inventory for variant %d: %v", productVariantID, err)
	}

	if stockStr == "" {
		return 0, fmt.Errorf("Inventory for variant %d not found in Redis", productVariantID)
	}

	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid inventory data format for variant %d: %v", productVariantID, err)
	}

	return stock, nil
}
