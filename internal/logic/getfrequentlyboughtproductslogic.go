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

type GetFrequentlyBoughtProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFrequentlyBoughtProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFrequentlyBoughtProductsLogic {
	return &GetFrequentlyBoughtProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFrequentlyBoughtProductsLogic) GetFrequentlyBoughtProducts(in *dropshipbe.GetFrequentlyBoughtProductsRequest) (*dropshipbe.ProductListResponse, error) {

	frequentlyBoughts, err := l.svcCtx.EcommerceRepo.GetFrequentlyBoughtProducts(l.ctx, in)
	if err != nil {
		l.Logger.Errorf("Error ProductID %d: %v", in.ProductId, err)
		return nil, err
	}
	var productItems []*dropshipbe.Product
	for _, p := range frequentlyBoughts {
		var flashSaleEndTime string
		var badge, saleLabel, saleTag string
		if p.BoughtWithProduct.Badge != nil {
			badge = *p.BoughtWithProduct.Badge
		} else {
			badge = ""
		}
		if p.BoughtWithProduct.SaleLabel != nil {
			saleLabel = *p.BoughtWithProduct.SaleLabel
		} else {
			saleLabel = ""
		}
		if p.BoughtWithProduct.SaleTag != nil {
			saleTag = *p.BoughtWithProduct.SaleTag
		} else {
			saleTag = ""
		}
		if p.BoughtWithProduct.FlashSaleEndTime != nil {
			flashSaleEndTime = p.BoughtWithProduct.FlashSaleEndTime.Format(time.RFC3339)
		} else {
			flashSaleEndTime = ""
		}
		var inStock bool
		if p.BoughtWithProduct.Status == "active" {
			inStock = true
		}
		productItems = append(productItems, &dropshipbe.Product{
			Id:          p.ID,
			CountryCode: p.BoughtWithProduct.Country.Code,
			Name:        p.BoughtWithProduct.Name,
			Slug:        p.BoughtWithProduct.Slug,
			WowDelay:    "",
			Metadata: map[string]string{
				"metadata": p.BoughtWithProduct.Metadata.String(),
			},
			Description:       p.BoughtWithProduct.Description,
			Rating:            float32(p.BoughtWithProduct.Rating),
			ReviewCount:       int32(p.BoughtWithProduct.ReviewCount),
			IsFeatured:        p.BoughtWithProduct.IsFeatured,
			IsTrending:        p.BoughtWithProduct.IsTrending,
			IsNew:             p.BoughtWithProduct.IsNew,
			Price:             float32(p.BoughtWithProduct.Price),
			InStock:           inStock,
			Status:            p.BoughtWithProduct.Status,
			Categories:        l.convertCategories(p.BoughtWithProduct.Categories),
			Galleries:         l.convertGaleries(p.BoughtWithProduct.Images),
			ProductPriceTiers: l.convertPriceTiers(p.BoughtWithProduct.PriceTiers),
			DescriptionImages: l.convertGaleries(p.BoughtWithProduct.Images),
			Options:           l.convertOptions(p.BoughtWithProduct.Options),
			Variants:          l.convertVariants(p.BoughtWithProduct.Variants),
			MetaTitle:         p.BoughtWithProduct.MetaTitle,
			MetaDescription:   p.BoughtWithProduct.MetaDescription,
			Vendor:            p.BoughtWithProduct.Vendor,
			ProductType:       p.BoughtWithProduct.ProductType,
			Badge:             badge,
			SaleLabel:         saleLabel,
			SaleTag:           saleTag,
			FlashSaleEndTime:  flashSaleEndTime,
			Sold:              int32(p.BoughtWithProduct.Sold),
			Tags:              l.convertTags(p.BoughtWithProduct.Tags),
			QuantityEnabled:   p.BoughtWithProduct.QuantityEnabled,
			QuickShop:         p.BoughtWithProduct.QuickShop,
			CreatedAt:         p.BoughtWithProduct.CreatedAt.Format(time.RFC3339),
		})
	}
	return &dropshipbe.ProductListResponse{
		Products: productItems,
	}, nil
}

func (l *GetFrequentlyBoughtProductsLogic) convertCategories(categories []model.Category) []*dropshipbe.Category {
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

func (l *GetFrequentlyBoughtProductsLogic) convertGaleries(images []model.ProductImage) []*dropshipbe.Gallery {
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
			l.Logger.Errorf("Error creating presigned URL for image %s: %v", i.ImageURL, err)
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

func (l *GetFrequentlyBoughtProductsLogic) convertPriceTiers(priceTiers []model.ProductPriceTier) []*dropshipbe.PriceTier {
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

func (l *GetFrequentlyBoughtProductsLogic) convertOptions(options []model.Option) []*dropshipbe.Option {
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
			ProductId:    o.ProductID,
			Name:         o.Name,
			Code:         o.Code,
			OptionValues: optionValueItems,
		})
	}
	return optionItems
}

func (l *GetFrequentlyBoughtProductsLogic) convertVariants(variants []model.Variant) []*dropshipbe.Variant {
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
			l.Logger.Errorf("Error creating presigned URL for variant image %s: %v", v.ImageURL, err)
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

func (l *GetFrequentlyBoughtProductsLogic) convertTags(jsonData datatypes.JSON) []string {
	var tags []string
	err := json.Unmarshal(jsonData, &tags)
	if err != nil {
		l.Logger.Errorf("Error converting tags: %v", err)
		return []string{}
	}
	return tags
}
