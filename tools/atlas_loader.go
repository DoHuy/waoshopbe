// tools/atlas_loader.go
package main

import (
	model "dropshipbe/model/schema"
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&model.User{},
		&model.FileMetadata{},
		&model.Country{},
		&model.Product{},
		&model.Option{},
		&model.OptionValue{},
		&model.Variant{},
		&model.Supplier{},
		&model.ProductMapping{},
		&model.Order{},
		&model.PurchaseOrder{},
		&model.PurchaseOrderItem{},
		&model.Shipment{},
		&model.Campaign{},
		&model.Coupon{},
		&model.CouponUsage{},
		&model.Transaction{},
		&model.Banner{},
		&model.Category{},
		&model.ProductImage{},
		&model.Policy{},
		&model.BlogCategory{},
		&model.BlogPost{},
		&model.ProductFAQ{},
		&model.ProductReview{},
		&model.ProductPriceTier{},
		&model.Slider{},
		&model.OrderItem{},
		&model.FrequentlyBought{},
		&model.Message{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(stmts)
}
