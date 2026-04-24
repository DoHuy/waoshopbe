package repository

import (
	"context"
	"fmt"
	"strings"

	"dropshipbe/common/constant"
	"dropshipbe/dropshipbe"
	model "dropshipbe/model/schema"

	"github.com/pgvector/pgvector-go"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EcommerceRepository interface {
	GetProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error)
	GetShop(ctx context.Context, request *dropshipbe.ShopSearchParams) ([]model.Product, error)
	GetFeaturedProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error)
	GetNewProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error)
	GetRelatedProducts(ctx context.Context, request *dropshipbe.GetRelatedProductsRequest) ([]model.Product, error)
	GetBannerItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Banner, error)
	GetBlogBySlug(ctx context.Context, request *dropshipbe.GetBlogBySlugRequest) (*model.BlogPost, error)
	GetBlogItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.BlogPost, error)
	GetCategoryItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Category, error)
	GetProductBySlug(ctx context.Context, request *dropshipbe.GetProductBySlugRequest) (*model.Product, error)
	GetProductByID(ctx context.Context, Id uint64) (*model.Product, error)
	GetProductFaqs(ctx context.Context, request *dropshipbe.GetProductFaqsRequest) ([]model.ProductFAQ, error)
	GetProductReviews(ctx context.Context, request *dropshipbe.GetProductReviewsRequest) ([]model.ProductReview, error)
	GetProductsByCategory(ctx context.Context, request *dropshipbe.GetProductsByCategoryRequest) ([]model.Product, error)
	GetSliderItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Slider, error)
	GetSocialProductVideos(ctx context.Context, request *dropshipbe.GetSocialProductVideoRequest) ([]model.ProductImage, error)
	GetVideoBanner(ctx context.Context, request *dropshipbe.DefaultRequest) (*model.Banner, error)
	CreateProductReview(ctx context.Context, request *dropshipbe.CreateProductReviewRequest) (*model.ProductReview, error)
	GetFrequentlyBoughtProducts(ctx context.Context, request *dropshipbe.GetFrequentlyBoughtProductsRequest) ([]model.FrequentlyBought, error)
	SaveMessage(ctx context.Context, ask, answer string) error
	TrackShipment(ctx context.Context, identifier string) (string, error)
	CheckVariantStock(ctx context.Context, productName string) (string, error)
	SearchProductsSemantic(ctx context.Context, embedding []float32, query string) (string, error)
	GetVariantsByIDs(ctx context.Context, ids []uint64) ([]model.Variant, error)
	CreateOrder(ctx context.Context, order *model.Order, orderItems []model.OrderItem) error
	UpdateOrderStatus(ctx context.Context, orderID uint64, newStatus string) error
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetTransactionByReference(ctx context.Context, reference string) (*model.Transaction, error)
	UpdateTransactionAndOrderStatus(ctx context.Context, transaction *model.Transaction, newTransactionStatus, newOrderStatus string, rawJSON datatypes.JSON, orderRawJSon map[string]interface{}) error
	GetProductByName(ctx context.Context, name string) (*model.Product, error)
	GetOrderDetailByOrderNumber(ctx context.Context, orderNumber string) (*model.Order, error)
	SaveChatbotInteraction(ctx context.Context, userMessage, botResponse string) error
	UpdateProductEmbedding(ctx context.Context, productID uint64, embedding []float32) error
}

type defaultEcommerceRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// UpdateProductEmbedding implements [EcommerceRepository].
func (d *defaultEcommerceRepository) UpdateProductEmbedding(ctx context.Context, productID uint64, embedding []float32) error {
	if len(embedding) == 0 {
		return fmt.Errorf("error updating product embedding for product %d: vector embedding is empty", productID)
	}

	vec := pgvector.NewVector(embedding)

	err := d.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", productID).
		Update("embedding", vec).Error

	if err != nil {
		return fmt.Errorf("error updating product embedding for product %d: %w", productID, err)
	}

	return nil
}

// SaveChatbotInteraction implements [EcommerceRepository].
func (d *defaultEcommerceRepository) SaveChatbotInteraction(ctx context.Context, userMessage string, botResponse string) error {
	interaction := &model.Message{
		Ask:    userMessage,
		Answer: botResponse,
	}
	return d.db.WithContext(ctx).Create(interaction).Error

}

// GetOrderDetailByOrderNumber implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetOrderDetailByOrderNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	var order model.Order
	err := d.db.WithContext(ctx).Preload("OrderItems").Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetProductByName implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetProductByName(ctx context.Context, name string) (*model.Product, error) {
	var product model.Product
	err := d.db.WithContext(ctx).Preload("Variants").Where("name ILIKE ?", "%"+name+"%").First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateTransactionAndOrderStatus implements [EcommerceRepository].
func (d *defaultEcommerceRepository) UpdateTransactionAndOrderStatus(ctx context.Context, transaction *model.Transaction, newTransactionStatus, newOrderStatus string, rawJSON datatypes.JSON, orderRawJson map[string]interface{}) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Model(&transaction).Updates(map[string]interface{}{
			"status":       newTransactionStatus,
			"raw_response": rawJSON,
		})

		tx.Model(transaction.Order).WithContext(ctx).Update("financial_status", newOrderStatus)

		if orderRawJson != nil {
			tx.Model(transaction.Order).WithContext(ctx).Updates(orderRawJson)
		}
		return nil
	})
}

// GetTransactionByReference implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetTransactionByReference(ctx context.Context, reference string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := d.db.WithContext(ctx).Preload("Order").Preload("Order.OrderItems").Where("transaction_reference = ?", reference).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// CreateTransaction implements [EcommerceRepository].
func (d *defaultEcommerceRepository) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	return d.db.WithContext(ctx).Create(transaction).Error
}

// UpdateOrderStatus implements [EcommerceRepository].
func (d *defaultEcommerceRepository) UpdateOrderStatus(ctx context.Context, orderID uint64, newStatus string) error {
	return d.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", orderID).Update("financial_status", newStatus).Error
}

// CreateOrder implements [EcommerceRepository].
func (d *defaultEcommerceRepository) CreateOrder(ctx context.Context, newOrder *model.Order, orderItems []model.OrderItem) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newOrder).Error; err != nil {
			return err
		}
		for i := range orderItems {
			orderItems[i].OrderID = newOrder.ID
		}

		if len(orderItems) > 0 {
			if err := tx.Create(&orderItems).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetVariantsByIDs implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetVariantsByIDs(ctx context.Context, ids []uint64) ([]model.Variant, error) {
	var variants []model.Variant
	if err := d.db.WithContext(ctx).Table("variants").Preload("Product").Where("id IN ?", ids).Find(&variants).Error; err != nil {
		return nil, fmt.Errorf("error fetching variants: %v", err)
	}
	return variants, nil
}

// CheckVariantStock implements [EcommerceRepository].
func (d *defaultEcommerceRepository) CheckVariantStock(ctx context.Context, productName string) (string, error) {
	panic("unimplemented")
}

// SaveMessage implements [EcommerceRepository].
func (d *defaultEcommerceRepository) SaveMessage(ctx context.Context, ask string, answer string) error {
	panic("unimplemented")
}

func (d *defaultEcommerceRepository) SearchProductsSemantic(ctx context.Context, embedding []float32, query string) (string, error) {
	var products []model.Product

	err := d.db.WithContext(ctx).
		Where("status = ?", "active").
		Order(clause.Expr{SQL: "embedding <=> ?", Vars: []interface{}{pgvector.NewVector(embedding)}}).
		Limit(5).
		Find(&products).Error

	if err != nil {
		return "Error occurred while searching for products.", err
	}

	if len(products) == 0 {
		return "No matching products found.", nil
	}

	var result strings.Builder
	result.WriteString("Here are the best matching products from the database:\n")
	for _, p := range products {
		fmt.Fprintf(&result, "- Name: %s | Price: %.2f | Description: %s\n", p.Name, p.Price, p.Description)
	}
	return result.String(), nil
}

// TrackShipment implements [EcommerceRepository].
func (d *defaultEcommerceRepository) TrackShipment(ctx context.Context, identifier string) (string, error) {
	var shipment model.Shipment
	err := d.db.WithContext(ctx).Table("shipments").Joins("JOIN orders ON orders.id = shipments.order_id").
		Where("orders.order_number = ? OR shipments.tracking_number = ?", identifier, identifier).
		First(&shipment).Error

	if err != nil {
		return fmt.Sprintf("No shipping information found for ID: %s", identifier), err
	}

	eta := "TBD"
	if shipment.EstimatedDeliveryDate != nil {
		eta = shipment.EstimatedDeliveryDate.Format("02/01/2006")
	}

	return fmt.Sprintf("Carrier: %s\nTracking Number: %s\nStatus: %s\nEstimated Delivery: %s",
		shipment.CarrierCode, shipment.TrackingNumber, shipment.Status, eta), nil
}

// GetFrequentlyBoughtProducts implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetFrequentlyBoughtProducts(ctx context.Context, request *dropshipbe.GetFrequentlyBoughtProductsRequest) ([]model.FrequentlyBought, error) {
	var items []model.FrequentlyBought

	cacheKey := constant.FrequentlyBoughtListKey(request.ProductId, request.CountryCode)

	err := d.cache.TakeCtx(ctx, &items, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.FrequentlyBought{}).
			Where("product_id = ? AND is_active = ?", request.ProductId, true).
			Order("sort_order ASC").
			Preload("BoughtWithProduct").
			Preload("BoughtWithProduct.Country").
			Preload("BoughtWithProduct.Categories").
			Preload("BoughtWithProduct.Images").
			Preload("BoughtWithProduct.PriceTiers").
			Preload("BoughtWithProduct.Options.OptionValues").
			Preload("BoughtWithProduct.Variants.OptionValues")

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetProductByID implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetProductByID(ctx context.Context, Id uint64) (*model.Product, error) {
	var product model.Product

	query := d.db.WithContext(ctx).Model(&model.Product{}).
		Preload("Country").
		Preload("Categories").
		Preload("Images").
		Preload("PriceTiers").
		Preload("Options.OptionValues").
		Preload("Variants.OptionValues")

	if err := query.Where("id = ?", Id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// CreateProductReview implements [EcommerceRepository].
func (d *defaultEcommerceRepository) CreateProductReview(ctx context.Context, request *dropshipbe.CreateProductReviewRequest) (*model.ProductReview, error) {
	review := &model.ProductReview{
		ProductID:    request.ProductId,
		AuthorName:   request.Name,
		AuthorEmail:  request.Email,
		AuthorAvatar: request.Avatar,
		Rating:       int(request.Rating),
		Content:      request.Comment,
		IsVerified:   true,
		Media: &model.ReviewMedia{
			Images: request.Images,
			Videos: request.Videos,
		},
		Status: "active",
	}

	if err := d.db.WithContext(ctx).Create(review).Error; err != nil {
		return nil, err
	}

	return review, nil
}

// GetSocialProductVideos implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetSocialProductVideos(ctx context.Context, request *dropshipbe.GetSocialProductVideoRequest) ([]model.ProductImage, error) {
	var videos []model.ProductImage
	cacheKey := constant.SocialProductVideoListKey(request.Id, request.CountryCode)

	err := d.cache.TakeCtx(ctx, &videos, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.ProductImage{}).Where("product_id = ? AND media_type = ?", request.Id, "social_video")
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return videos, nil
}

// GetVideoBanner implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetVideoBanner(ctx context.Context, request *dropshipbe.DefaultRequest) (*model.Banner, error) {
	var banner *model.Banner
	cacheKey := constant.VideoBannerKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &banner, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.Banner{}).Where("is_active = ?", true)
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = banner.country_code").
				Where("countries.code = ?", request.CountryCode)
		}
		return query.First(v).Error
	})

	if err != nil {
		return nil, err
	}

	return banner, nil
}

func (d *defaultEcommerceRepository) GetSliderItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Slider, error) {
	var sliders []model.Slider
	cacheKey := constant.SliderItemListKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &sliders, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.Slider{}).Where("is_active = ?", true)
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = sliders.country_code").
				Where("countries.code = ?", request.CountryCode)
		}
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return sliders, nil
}

// GetShop implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetShop(ctx context.Context, request *dropshipbe.ShopSearchParams) ([]model.Product, error) {
	var products []model.Product
	cacheKey := constant.ShopSearchKey(request.IsFeatured, request.IsNew, request.IsOnSale, request.IsTrending, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ?", "active")

		if request.IsFeatured {
			query = query.Where("is_featured = ?", true)
		}
		if request.IsNew {
			query = query.Where("is_new = ?", true)
		}
		if request.IsOnSale {
			query = query.Where("is_on_sale = ?", true)
		}
		if request.IsTrending {
			query = query.Where("is_trending = ?", true)
		}
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return products, nil
}

// GetProductsByCategory implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetProductsByCategory(ctx context.Context, request *dropshipbe.GetProductsByCategoryRequest) ([]model.Product, error) {
	var products []model.Product
	cacheKey := constant.ProductListByCategoryKey(request.Category, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Joins("JOIN product_categories ON product_categories.product_id = products.id").
			Joins("JOIN categories ON categories.id = product_categories.category_id").
			Where("categories.slug = ? AND products.status = ?", request.Category, "active")

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return products, nil
}

// GetProductReviews implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetProductReviews(ctx context.Context, request *dropshipbe.GetProductReviewsRequest) ([]model.ProductReview, error) {
	var reviews []model.ProductReview
	cacheKey := constant.ProductReviewListKey(request.Id, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &reviews, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.ProductReview{}).Where("product_id = ? AND is_verified = ? AND status = ?", request.Id, true, "active")
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (d *defaultEcommerceRepository) GetProductFaqs(ctx context.Context, request *dropshipbe.GetProductFaqsRequest) ([]model.ProductFAQ, error) {
	var faqs []model.ProductFAQ
	cacheKey := constant.ProductFaqListKey(request.Id, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &faqs, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.ProductFAQ{}).Where("product_id = ?", request.Id)
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return faqs, nil
}

// GetProductBySlug implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetProductBySlug(ctx context.Context, request *dropshipbe.GetProductBySlugRequest) (*model.Product, error) {
	var product model.Product
	cacheKey := constant.ProductDetailKey(request.Slug, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &product, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues")

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Where("slug = ?", request.Slug).First(v).Error
	})

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetCategoryItems implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetCategoryItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Category, error) {
	var categories []model.Category
	cacheKey := constant.CategoryListKey(request.CountryCode)
	err := d.cache.TakeCtx(ctx, &categories, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.Category{})
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = categories.country_code").
				Where("countries.code = ?", request.CountryCode)
		}
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetBlogItems implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetBlogItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.BlogPost, error) {
	var blogPosts []model.BlogPost
	cacheKey := constant.BlogPostListByCountryKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &blogPosts, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.BlogPost{}).
			Preload("Category").
			Preload("Country").
			Where("published_at IS NOT NULL")

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = blog_posts.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return blogPosts, nil
}

func (d *defaultEcommerceRepository) GetBlogBySlug(ctx context.Context, request *dropshipbe.GetBlogBySlugRequest) (*model.BlogPost, error) {
	var blogPost model.BlogPost
	cacheKey := constant.BlogPostBySlugKey(request.Slug, request.CountryCode)

	err := d.cache.TakeCtx(ctx, &blogPost, cacheKey, func(v any) error {
		query := d.db.WithContext(ctx).Model(&model.BlogPost{})
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = blog_posts.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.
			Preload("Category").
			Preload("Country").
			Where("slug = ?", request.Slug).First(v).Error
	})

	if err != nil {
		return nil, err
	}

	return &blogPost, nil
}

func (d *defaultEcommerceRepository) GetBannerItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Banner, error) {
	var banners []model.Banner
	cacheKey := constant.BannerItemListKey()

	err := d.cache.TakeCtx(ctx, &banners, cacheKey, func(v any) error {
		return d.db.WithContext(ctx).Model(&model.Banner{}).Where("is_active = ?", true).
			Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return banners, nil
}

func NewEcommerceRepository(db *gorm.DB, c cache.Cache) EcommerceRepository {
	return &defaultEcommerceRepository{
		db:    db,
		cache: c,
	}
}

func (d *defaultEcommerceRepository) GetProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error) {
	var products []model.Product

	cacheKey := constant.ProductListByCountryKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {

		query := d.db.WithContext(ctx).Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ?", "active")

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (d *defaultEcommerceRepository) GetFeaturedProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error) {
	var products []model.Product

	cacheKey := constant.FeaturedProductListKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {

		query := d.db.WithContext(ctx).Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ? AND is_featured = ?", "active", true)

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (d *defaultEcommerceRepository) GetRelatedProducts(ctx context.Context, request *dropshipbe.GetRelatedProductsRequest) ([]model.Product, error) {
	var products []model.Product
	fmt.Println("id:", request.Id)
	cacheKey := constant.RelatedProductListKey(request.Id, request.CountryCode)

	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {

		query := d.db.WithContext(ctx).Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ? AND related_product_id = ?", "active", request.Id)

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return products, nil
}
func (d *defaultEcommerceRepository) GetNewProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error) {
	var products []model.Product

	cacheKey := constant.NewProductListKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {

		query := d.db.WithContext(ctx).Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ? AND is_new = ?", "active", true)

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return products, nil
}
