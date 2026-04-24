package logic

import (
	"context"
	"encoding/json"
	"time"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"
	model "dropshipbe/model/schema"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
)

type GetProductsByCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductsByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductsByCategoryLogic {
	return &GetProductsByCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Products ---
func (l *GetProductsByCategoryLogic) convertCategories(categories []model.Category) []*dropshipbe.Category {
	var categoryItems []*dropshipbe.Category
	for _, c := range categories {
		categoryItems = append(categoryItems, &dropshipbe.Category{
			Id:          c.ID,
			Name:        c.Name,
			CountryCode: c.CountryCode,
			Slug:        c.Slug,
			Description: c.Description,
			ImageUrl:    "",
			IsActive:    *c.IsActive,
			ButtonText:  "",
			Alt:         "",
		})
	}
	return categoryItems
}

func (l *GetProductsByCategoryLogic) convertGaleries(images []model.ProductImage) []*dropshipbe.Gallery {
	var imageItems []*dropshipbe.Gallery
	expirationDuration := time.Duration(l.svcCtx.Config.R2.LinkExpiration) * time.Minute
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, i := range images {

		presignedReq, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(i.ImageURL),
		}, s3.WithPresignExpires(expirationDuration))

		if err != nil {
			l.Logger.Errorf("Error  presigned URL image %s: %v", i.ImageURL, err)
			continue
		}

		presignedVideoReq, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(i.VideoURL),
		}, s3.WithPresignExpires(expirationDuration))

		if err != nil {
			l.Logger.Errorf("Error creating presigned URL for video %s: %v", i.VideoURL, err)
			continue
		}

		imageItems = append(imageItems, &dropshipbe.Gallery{
			Id:        i.ID,
			ProductId: i.ProductID,
			ImageUrl:  presignedReq.URL,
			VideoUrl:  presignedVideoReq.URL,
			AltText:   i.AltText,
			Position:  int32(i.Position),
		})

	}
	return imageItems
}

func (l *GetProductsByCategoryLogic) convertPriceTiers(priceTiers []model.ProductPriceTier) []*dropshipbe.PriceTier {
	var priceTierItems []*dropshipbe.PriceTier
	for _, pt := range priceTiers {
		priceTierItems = append(priceTierItems, &dropshipbe.PriceTier{
			Id:        pt.ID,
			ProductId: pt.ProductID,
			Price:     float32(pt.Price),
			Savings:   pt.SavingsText,
			Qty:       int32(pt.Qty),
			Tag:       pt.Tag,
			TagClass:  pt.TagClass,
			CreatedAt: pt.CreatedAt.Format(time.RFC3339),
		})
	}
	return priceTierItems
}

func (l *GetProductsByCategoryLogic) convertOptions(options []model.Option) []*dropshipbe.Option {
	var optionItems []*dropshipbe.Option
	for _, o := range options {
		var optionValueItems []*dropshipbe.OptionValue
		for _, ov := range o.OptionValues {
			optionValueItems = append(optionValueItems, &dropshipbe.OptionValue{
				Id:        ov.ID,
				Value:     ov.Value,
				ColorCode: ov.ColorCode,
				OptionId:  ov.OptionID,
			})
		}
		optionItems = append(optionItems, &dropshipbe.Option{
			Id:           o.ID,
			Name:         o.Name,
			Code:         o.Code,
			OptionValues: optionValueItems,
		})
	}
	return optionItems
}

func (l *GetProductsByCategoryLogic) convertVariants(variants []model.Variant) []*dropshipbe.Variant {
	var variantItems []*dropshipbe.Variant

	expirationDuration := time.Duration(l.svcCtx.Config.R2.LinkExpiration) * time.Minute
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	for _, v := range variants {
		var variantOptionValueItems []*dropshipbe.VariantOption
		for _, ov := range v.OptionValues {
			variantOptionValueItems = append(variantOptionValueItems, &dropshipbe.VariantOption{
				OptionId:      ov.OptionID,
				OptionCode:    ov.ColorCode,
				OptionValueId: ov.ID,
				OptionValue:   ov.Value,
			})
		}

		presignedImage, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(v.ImageURL),
		}, s3.WithPresignExpires(expirationDuration))

		if err != nil {
			l.Logger.Errorf("Error creating presigned URL for image variant %s: %v", v.ImageURL, err)
			continue
		}
		variantItems = append(variantItems, &dropshipbe.Variant{
			Id:             v.ID,
			Sku:            v.Sku,
			ProductId:      v.ProductID,
			ImageUrl:       presignedImage.URL,
			Price:          float32(v.Price),
			Barcode:        v.Barcode,
			CompareAtPrice: float32(v.CompareAtPrice),
			CostPrice:      float32(v.CostPrice),
			StockQuantity:  int32(v.StockQuantity),
			Options:        variantOptionValueItems,
			IsActive:       *v.IsActive,
			CreatedAt:      v.CreatedAt.Format(time.RFC3339),
		})
	}
	return variantItems
}

func (l *GetProductsByCategoryLogic) convertTags(jsonData datatypes.JSON) []string {
	var tags []string
	err := json.Unmarshal(jsonData, &tags)
	if err != nil {
		l.Logger.Errorf("Error converting tags: %v", err)
		return []string{}
	}
	return tags
}

func (l *GetProductsByCategoryLogic) GetProductsByCategory(in *dropshipbe.GetProductsByCategoryRequest) (*dropshipbe.ProductListResponse, error) {
	// todo: add your logic here and delete this line
	products, err := l.svcCtx.EcommerceRepo.GetProductsByCategory(l.ctx, in)
	if err != nil {
		l.Logger.Errorf("Error fetching products: %v", err)
		return nil, err
	}

	var productItems []*dropshipbe.Product
	for _, p := range products {
		var flashSaleEndTime string
		var badge, saleLabel, saleTag string
		if p.Badge != nil {
			badge = *p.Badge
		} else {
			badge = ""
		}
		if p.SaleLabel != nil {
			saleLabel = *p.SaleLabel
		} else {
			saleLabel = ""
		}
		if p.SaleTag != nil {
			saleTag = *p.SaleTag
		} else {
			saleTag = ""
		}
		if p.FlashSaleEndTime != nil {
			flashSaleEndTime = p.FlashSaleEndTime.Format(time.RFC3339)
		} else {
			flashSaleEndTime = ""
		}
		productItems = append(productItems, &dropshipbe.Product{
			Id:          p.ID,
			CountryCode: p.Country.Code,
			Name:        p.Name,
			Slug:        p.Slug,
			WowDelay:    "",
			Metadata: map[string]string{
				"metadata": p.Metadata.String(),
			},
			Description: p.Description,
			Rating:      float32(p.Rating),
			ReviewCount: int32(p.ReviewCount),
			IsFeatured:  p.IsFeatured,
			IsTrending:  p.IsTrending,
			IsNew:       p.IsNew,
			Price:       float32(p.Price),

			Status:            p.Status,
			Categories:        l.convertCategories(p.Categories),
			Galleries:         l.convertGaleries(p.Images),
			ProductPriceTiers: l.convertPriceTiers(p.PriceTiers),
			DescriptionImages: l.convertGaleries(p.Images),
			Options:           l.convertOptions(p.Options),
			Variants:          l.convertVariants(p.Variants),
			MetaTitle:         p.MetaTitle,
			MetaDescription:   p.MetaDescription,
			Vendor:            p.Vendor,
			ProductType:       p.ProductType,
			Badge:             badge,
			SaleLabel:         saleLabel,
			SaleTag:           saleTag,
			FlashSaleEndTime:  flashSaleEndTime,
			Sold:              int32(p.Sold),
			Tags:              l.convertTags(p.Tags),
			QuantityEnabled:   p.QuantityEnabled,
			QuickShop:         p.QuickShop,
			CreatedAt:         p.CreatedAt.Format(time.RFC3339),
		})
	}

	return &dropshipbe.ProductListResponse{
		Products: productItems,
	}, nil
}
