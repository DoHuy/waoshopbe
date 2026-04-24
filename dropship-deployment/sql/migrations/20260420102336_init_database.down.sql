-- reverse: create "variant_value_map" table
DROP TABLE "public"."variant_value_map";
-- reverse: create index "idx_transactions_deleted_at" to table: "transactions"
DROP INDEX "public"."idx_transactions_deleted_at";
-- reverse: create "transactions" table
DROP TABLE "public"."transactions";
-- reverse: create index "idx_sliders_deleted_at" to table: "sliders"
DROP INDEX "public"."idx_sliders_deleted_at";
-- reverse: create index "idx_sliders_country_code" to table: "sliders"
DROP INDEX "public"."idx_sliders_country_code";
-- reverse: create "sliders" table
DROP TABLE "public"."sliders";
-- reverse: create index "idx_shipments_deleted_at" to table: "shipments"
DROP INDEX "public"."idx_shipments_deleted_at";
-- reverse: create "shipments" table
DROP TABLE "public"."shipments";
-- reverse: create index "idx_purchase_order_items_deleted_at" to table: "purchase_order_items"
DROP INDEX "public"."idx_purchase_order_items_deleted_at";
-- reverse: create "purchase_order_items" table
DROP TABLE "public"."purchase_order_items";
-- reverse: create index "idx_purchase_orders_deleted_at" to table: "purchase_orders"
DROP INDEX "public"."idx_purchase_orders_deleted_at";
-- reverse: create "purchase_orders" table
DROP TABLE "public"."purchase_orders";
-- reverse: create index "idx_product_reviews_status" to table: "product_reviews"
DROP INDEX "public"."idx_product_reviews_status";
-- reverse: create index "idx_product_reviews_rating" to table: "product_reviews"
DROP INDEX "public"."idx_product_reviews_rating";
-- reverse: create index "idx_product_reviews_product_id" to table: "product_reviews"
DROP INDEX "public"."idx_product_reviews_product_id";
-- reverse: create index "idx_product_reviews_deleted_at" to table: "product_reviews"
DROP INDEX "public"."idx_product_reviews_deleted_at";
-- reverse: create index "idx_product_reviews_author_email" to table: "product_reviews"
DROP INDEX "public"."idx_product_reviews_author_email";
-- reverse: create "product_reviews" table
DROP TABLE "public"."product_reviews";
-- reverse: create index "idx_product_price_tiers_product_id" to table: "product_price_tiers"
DROP INDEX "public"."idx_product_price_tiers_product_id";
-- reverse: create index "idx_product_price_tiers_deleted_at" to table: "product_price_tiers"
DROP INDEX "public"."idx_product_price_tiers_deleted_at";
-- reverse: create "product_price_tiers" table
DROP TABLE "public"."product_price_tiers";
-- reverse: create index "uniq_mapping_local_var" to table: "product_mappings"
DROP INDEX "public"."uniq_mapping_local_var";
-- reverse: create index "idx_product_mappings_deleted_at" to table: "product_mappings"
DROP INDEX "public"."idx_product_mappings_deleted_at";
-- reverse: create "product_mappings" table
DROP TABLE "public"."product_mappings";
-- reverse: create index "idx_suppliers_deleted_at" to table: "suppliers"
DROP INDEX "public"."idx_suppliers_deleted_at";
-- reverse: create "suppliers" table
DROP TABLE "public"."suppliers";
-- reverse: create "product_images" table
DROP TABLE "public"."product_images";
-- reverse: create index "idx_product_faqs_product_id" to table: "product_faqs"
DROP INDEX "public"."idx_product_faqs_product_id";
-- reverse: create "product_faqs" table
DROP TABLE "public"."product_faqs";
-- reverse: create "product_categories" table
DROP TABLE "public"."product_categories";
-- reverse: create index "idx_policy_country_type" to table: "policies"
DROP INDEX "public"."idx_policy_country_type";
-- reverse: create "policies" table
DROP TABLE "public"."policies";
-- reverse: create index "idx_order_items_variant_id" to table: "order_items"
DROP INDEX "public"."idx_order_items_variant_id";
-- reverse: create index "idx_order_items_product_id" to table: "order_items"
DROP INDEX "public"."idx_order_items_product_id";
-- reverse: create index "idx_order_items_order_id" to table: "order_items"
DROP INDEX "public"."idx_order_items_order_id";
-- reverse: create index "idx_order_items_deleted_at" to table: "order_items"
DROP INDEX "public"."idx_order_items_deleted_at";
-- reverse: create "order_items" table
DROP TABLE "public"."order_items";
-- reverse: create index "idx_option_values_option_id" to table: "option_values"
DROP INDEX "public"."idx_option_values_option_id";
-- reverse: create "option_values" table
DROP TABLE "public"."option_values";
-- reverse: create index "idx_options_product_id" to table: "options"
DROP INDEX "public"."idx_options_product_id";
-- reverse: create "options" table
DROP TABLE "public"."options";
-- reverse: create index "idx_product_bought_with" to table: "frequently_boughts"
DROP INDEX "public"."idx_product_bought_with";
-- reverse: create "frequently_boughts" table
DROP TABLE "public"."frequently_boughts";
-- reverse: create index "idx_coupon_usages_deleted_at" to table: "coupon_usages"
DROP INDEX "public"."idx_coupon_usages_deleted_at";
-- reverse: create "coupon_usages" table
DROP TABLE "public"."coupon_usages";
-- reverse: create index "idx_orders_number" to table: "orders"
DROP INDEX "public"."idx_orders_number";
-- reverse: create index "idx_orders_deleted_at" to table: "orders"
DROP INDEX "public"."idx_orders_deleted_at";
-- reverse: create "orders" table
DROP TABLE "public"."orders";
-- reverse: create "coupon_items" table
DROP TABLE "public"."coupon_items";
-- reverse: create index "idx_variants_sku" to table: "variants"
DROP INDEX "public"."idx_variants_sku";
-- reverse: create index "idx_variants_product_id" to table: "variants"
DROP INDEX "public"."idx_variants_product_id";
-- reverse: create "variants" table
DROP TABLE "public"."variants";
-- reverse: create index "idx_products_status" to table: "products"
DROP INDEX "public"."idx_products_status";
-- reverse: create index "idx_products_related_product_id" to table: "products"
DROP INDEX "public"."idx_products_related_product_id";
-- reverse: create index "idx_product_country_slug" to table: "products"
DROP INDEX "public"."idx_product_country_slug";
-- reverse: create "products" table
DROP TABLE "public"."products";
-- reverse: create index "idx_coupons_deleted_at" to table: "coupons"
DROP INDEX "public"."idx_coupons_deleted_at";
-- reverse: create index "idx_coupons_code" to table: "coupons"
DROP INDEX "public"."idx_coupons_code";
-- reverse: create "coupons" table
DROP TABLE "public"."coupons";
-- reverse: create index "idx_campaigns_deleted_at" to table: "campaigns"
DROP INDEX "public"."idx_campaigns_deleted_at";
-- reverse: create "campaigns" table
DROP TABLE "public"."campaigns";
-- reverse: create index "idx_category_active" to table: "categories"
DROP INDEX "public"."idx_category_active";
-- reverse: create index "idx_categories_parent_id" to table: "categories"
DROP INDEX "public"."idx_categories_parent_id";
-- reverse: create index "idx_cat_country_slug" to table: "categories"
DROP INDEX "public"."idx_cat_country_slug";
-- reverse: create "categories" table
DROP TABLE "public"."categories";
-- reverse: create index "idx_blog_published" to table: "blog_posts"
DROP INDEX "public"."idx_blog_published";
-- reverse: create index "idx_blog_posts_slug" to table: "blog_posts"
DROP INDEX "public"."idx_blog_posts_slug";
-- reverse: create index "idx_blog_posts_deleted_at" to table: "blog_posts"
DROP INDEX "public"."idx_blog_posts_deleted_at";
-- reverse: create index "idx_blog_posts_country_code" to table: "blog_posts"
DROP INDEX "public"."idx_blog_posts_country_code";
-- reverse: create index "idx_blog_posts_category_id" to table: "blog_posts"
DROP INDEX "public"."idx_blog_posts_category_id";
-- reverse: create "blog_posts" table
DROP TABLE "public"."blog_posts";
-- reverse: create index "idx_blog_categories_slug" to table: "blog_categories"
DROP INDEX "public"."idx_blog_categories_slug";
-- reverse: create index "idx_blog_categories_country_code" to table: "blog_categories"
DROP INDEX "public"."idx_blog_categories_country_code";
-- reverse: create "blog_categories" table
DROP TABLE "public"."blog_categories";
-- reverse: create index "idx_banners_deleted_at" to table: "banners"
DROP INDEX "public"."idx_banners_deleted_at";
-- reverse: create index "idx_banner_active" to table: "banners"
DROP INDEX "public"."idx_banner_active";
-- reverse: create "banners" table
DROP TABLE "public"."banners";
-- reverse: create "messages" table
DROP TABLE "public"."messages";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create index "idx_country_active" to table: "countries"
DROP INDEX "public"."idx_country_active";
-- reverse: create "countries" table
DROP TABLE "public"."countries";
-- reverse: create "file_metadata" table
DROP TABLE "public"."file_metadata";
