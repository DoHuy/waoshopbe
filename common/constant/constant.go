package constant

import (
	"fmt"
	"time"
)

const (
	ProjectPrefix = "dropship:"

	PrefixProduct  = ProjectPrefix + "product:"
	PrefixCategory = ProjectPrefix + "category:"
	PrefixBlog     = ProjectPrefix + "blog:"
)

func ProductListByCountryKey(countryCode string) string {
	if countryCode == "" {
		return PrefixProduct + "list:all"
	}
	return fmt.Sprintf("%slist:country:%s", PrefixProduct, countryCode)
}

func ProductDetailKey(slug string, country_code string) string {
	return fmt.Sprintf("%sdetail:slug:%s:country:%s", PrefixProduct, slug, country_code)
}

// CategoryListKey tạo key cho danh sách danh mục
func CategoryListKey(country_code string) string {
	if country_code == "" {
		return PrefixCategory + "list:all"
	}
	return fmt.Sprintf("%slist:country:%s", PrefixCategory, country_code)
}

func BannerItemListKey() string {
	return PrefixProduct + "banner:list:all"
}

func BlogPostBySlugKey(slug string, country_code string) string {
	return fmt.Sprintf("%s:country:%s:detail:slug:%s", PrefixBlog, country_code, slug)
}

func BlogPostListByCountryKey(countryCode string) string {
	if countryCode == "" {
		return PrefixBlog + "list:all"
	}
	return fmt.Sprintf("%slist:country:%s", PrefixBlog, countryCode)
}

func FeaturedProductListKey(countryCode string) string {
	if countryCode == "" {
		return PrefixProduct + "featured:list:all"
	}
	return fmt.Sprintf("%sfeatured:list:country:%s", PrefixProduct, countryCode)
}

func NewProductListKey(countryCode string) string {
	if countryCode == "" {
		return PrefixProduct + "new:list:all"
	}
	return fmt.Sprintf("%snew:list:country:%s", PrefixProduct, countryCode)
}

func RelatedProductListKey(relatedId uint64, countryCode string) string {
	if countryCode == "" {
		return fmt.Sprintf("%srelated:list:relatedId:%d:all", PrefixProduct, relatedId)
	}
	return fmt.Sprintf("%srelated:list:relatedId:%d:country:%s", PrefixProduct, relatedId, countryCode)
}

func ProductFaqListKey(productId uint64, countryCode string) string {
	return fmt.Sprintf("%sfaq:list:product:%d:country:%s", PrefixProduct, productId, countryCode)
}

func ProductReviewListKey(productId uint64, countryCode string) string {
	return fmt.Sprintf("%sreview:list:product:%d:country:%s", PrefixProduct, productId, countryCode)
}

func ProductListByCategoryKey(categoryId string, countryCode string) string {
	if countryCode == "" {
		return fmt.Sprintf("%slist:category:%s:all", PrefixProduct, categoryId)
	}
	return fmt.Sprintf("%slist:category:%s:country:%s", PrefixProduct, categoryId, countryCode)
}

func ShopSearchKey(isFeatured, isNew, isOnSale, isTrending bool, countryCode string) string {
	if countryCode == "" {
		return fmt.Sprintf("%sshop:search:featured:%t:new:%t:on_sale:%t:trending:%t:all", PrefixProduct, isFeatured, isNew, isOnSale, isTrending)
	}
	return fmt.Sprintf("%sshop:search:featured:%t:new:%t:on_sale:%t:trending:%t:country:%s", PrefixProduct, isFeatured, isNew, isOnSale, isTrending, countryCode)
}

func SliderItemListKey(country_code string) string {
	if country_code == "" {
		return PrefixProduct + "slider:list:all"
	}
	return fmt.Sprintf("%sslider:list:country:%s", PrefixProduct, country_code)
}

func SocialProductVideoListKey(productId uint64, country_code string) string {
	return fmt.Sprintf("%ssocial_video:list:product:%d:country:%s", PrefixProduct, productId, country_code)
}

func VideoBannerKey(country_code string) string {
	return fmt.Sprintf("%svideo_banner:country:%s", PrefixProduct, country_code)
}

func ProductVariantStock(variantID string) string {
	return fmt.Sprintf("%s:VariantID:stock:%s", PrefixProduct, variantID)
}

func CreateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().Unix())
}

func FrequentlyBoughtListKey(productId uint64, countryCode string) string {
	if countryCode == "" {
		return fmt.Sprintf("%sfrequentlyBought:list:relatedId:%d:all", PrefixProduct, productId)
	}
	return fmt.Sprintf("%sfrequentlyBought:list:relatedId:%d:country:%s", PrefixProduct, productId, countryCode)
}

func PaypalWebhookProcessedKey(eventID string) string {
	return fmt.Sprintf("paypal_webhook_processed:%s", eventID)
}

func OrderPaidLockKey(orderID string) string {
	return fmt.Sprintf("order_paid_lock:%s", orderID)
}
