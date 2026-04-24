-- create "file_metadata" table
CREATE TABLE "public"."file_metadata" (
  "id" bigserial NOT NULL,
  "filename" text NULL,
  "size" integer NULL,
  "content_type" text NULL,
  "file_key" text NULL,
  "uploaded_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "countries" table
CREATE TABLE "public"."countries" (
  "code" character(2) NOT NULL,
  "name" character varying(100) NOT NULL,
  "currency" character varying(3) NOT NULL,
  "language_code" character(2) NULL DEFAULT 'vi',
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("code")
);
-- create index "idx_country_active" to table: "countries"
CREATE INDEX "idx_country_active" ON "public"."countries" ("is_active") WHERE (is_active = true);
-- create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "role" text NULL,
  "revoke_tokens_before" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_username" UNIQUE ("username")
);
-- create "messages" table
CREATE TABLE "public"."messages" (
  "id" bigserial NOT NULL,
  "ask" text NOT NULL,
  "answer" text NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "banners" table
CREATE TABLE "public"."banners" (
  "id" bigserial NOT NULL,
  "title" character varying(255) NOT NULL,
  "country_code" character(2) NULL DEFAULT 'VN',
  "image_url" character varying(500) NOT NULL,
  "video_url" character varying(500) NULL,
  "alt" character varying(255) NULL,
  "description" text NULL,
  "video_type" character varying(50) NULL,
  "link_url" character varying(500) NULL,
  "position" character varying(50) NOT NULL,
  "sort_order" bigint NULL DEFAULT 0,
  "heading" character varying(255) NULL,
  "sub_heading" text NULL,
  "button_text" character varying(50) NULL,
  "is_active" boolean NULL DEFAULT true,
  "start_date" timestamptz NULL,
  "end_date" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_banners_country" FOREIGN KEY ("country_code") REFERENCES "public"."countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_banner_active" to table: "banners"
CREATE INDEX "idx_banner_active" ON "public"."banners" ("is_active") WHERE (is_active = true);
-- create index "idx_banners_deleted_at" to table: "banners"
CREATE INDEX "idx_banners_deleted_at" ON "public"."banners" ("deleted_at");
-- create "blog_categories" table
CREATE TABLE "public"."blog_categories" (
  "id" bigserial NOT NULL,
  "country_code" character(2) NOT NULL,
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "description" text NULL,
  "css_class" character varying(100) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_blog_categories_country" FOREIGN KEY ("country_code") REFERENCES "public"."countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_blog_categories_country_code" to table: "blog_categories"
CREATE INDEX "idx_blog_categories_country_code" ON "public"."blog_categories" ("country_code");
-- create index "idx_blog_categories_slug" to table: "blog_categories"
CREATE INDEX "idx_blog_categories_slug" ON "public"."blog_categories" ("slug");
-- create "blog_posts" table
CREATE TABLE "public"."blog_posts" (
  "id" bigserial NOT NULL,
  "country_code" character(2) NOT NULL,
  "category_id" bigint NULL,
  "title" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "excerpt" text NULL,
  "content" text NOT NULL,
  "meta_description" character varying(500) NULL,
  "image_url" character varying(500) NULL,
  "image_alt" character varying(255) NULL,
  "image_width" bigint NULL DEFAULT 0,
  "image_height" bigint NULL DEFAULT 0,
  "author_name" character varying(100) NULL,
  "author_avatar" character varying(500) NULL,
  "tags" jsonb NULL,
  "is_published" boolean NULL DEFAULT false,
  "published_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_blog_categories_blog_posts" FOREIGN KEY ("category_id") REFERENCES "public"."blog_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_blog_posts_country" FOREIGN KEY ("country_code") REFERENCES "public"."countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_blog_posts_category_id" to table: "blog_posts"
CREATE INDEX "idx_blog_posts_category_id" ON "public"."blog_posts" ("category_id");
-- create index "idx_blog_posts_country_code" to table: "blog_posts"
CREATE INDEX "idx_blog_posts_country_code" ON "public"."blog_posts" ("country_code");
-- create index "idx_blog_posts_deleted_at" to table: "blog_posts"
CREATE INDEX "idx_blog_posts_deleted_at" ON "public"."blog_posts" ("deleted_at");
-- create index "idx_blog_posts_slug" to table: "blog_posts"
CREATE INDEX "idx_blog_posts_slug" ON "public"."blog_posts" ("slug");
-- create index "idx_blog_published" to table: "blog_posts"
CREATE INDEX "idx_blog_published" ON "public"."blog_posts" ("is_published", "published_at");
-- create "categories" table
CREATE TABLE "public"."categories" (
  "id" bigserial NOT NULL,
  "parent_id" bigint NULL,
  "country_code" character(2) NOT NULL,
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "description" text NULL,
  "image_url" character varying(255) NULL,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_children" FOREIGN KEY ("parent_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_categories_country" FOREIGN KEY ("country_code") REFERENCES "public"."countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_cat_country_slug" to table: "categories"
CREATE INDEX "idx_cat_country_slug" ON "public"."categories" ("country_code", "slug");
-- create index "idx_categories_parent_id" to table: "categories"
CREATE INDEX "idx_categories_parent_id" ON "public"."categories" ("parent_id");
-- create index "idx_category_active" to table: "categories"
CREATE INDEX "idx_category_active" ON "public"."categories" ("is_active") WHERE (is_active = true);
-- create "campaigns" table
CREATE TABLE "public"."campaigns" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_campaigns_deleted_at" to table: "campaigns"
CREATE INDEX "idx_campaigns_deleted_at" ON "public"."campaigns" ("deleted_at");
-- create "coupons" table
CREATE TABLE "public"."coupons" (
  "id" bigserial NOT NULL,
  "campaign_id" bigint NOT NULL,
  "code" character varying(50) NOT NULL,
  "discount_type" character varying(30) NOT NULL,
  "value" numeric(15,2) NOT NULL,
  "min_order_value" numeric(15,2) NULL DEFAULT 0,
  "max_discount_amount" numeric(15,2) NULL,
  "target_type" character varying(30) NULL DEFAULT 'specific_items',
  "usage_limit" bigint NULL DEFAULT 0,
  "usage_limit_per_user" bigint NULL DEFAULT 1,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_campaigns_coupons" FOREIGN KEY ("campaign_id") REFERENCES "public"."campaigns" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_coupons_code" to table: "coupons"
CREATE UNIQUE INDEX "idx_coupons_code" ON "public"."coupons" ("code");
-- create index "idx_coupons_deleted_at" to table: "coupons"
CREATE INDEX "idx_coupons_deleted_at" ON "public"."coupons" ("deleted_at");
-- create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "country_code" character(2) NOT NULL,
  "related_product_id" bigint NULL,
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "metadata" jsonb NULL,
  "description" text NULL,
  "status" character varying(20) NULL DEFAULT 'draft',
  "price" numeric(15,2) NOT NULL,
  "is_featured" boolean NULL DEFAULT false,
  "is_new" boolean NULL DEFAULT false,
  "is_trending" boolean NULL DEFAULT false,
  "is_on_sale" boolean NULL DEFAULT false,
  "meta_title" character varying(255) NULL,
  "meta_description" character varying(500) NULL,
  "vendor" character varying(100) NULL,
  "product_type" character varying(100) NULL,
  "badge" character varying(50) NULL,
  "sale_label" character varying(50) NULL,
  "sale_tag" character varying(100) NULL,
  "flash_sale_end_time" timestamptz NULL,
  "sold" bigint NULL DEFAULT 0,
  "rating" numeric NULL DEFAULT 0,
  "review_count" bigint NULL DEFAULT 0,
  "tags" jsonb NULL,
  "quantity_enabled" boolean NULL DEFAULT true,
  "quick_shop" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "embedding" public.vector(1536) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_country" FOREIGN KEY ("country_code") REFERENCES "public"."countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_product_country_slug" to table: "products"
CREATE INDEX "idx_product_country_slug" ON "public"."products" ("country_code", "slug");
-- create index "idx_products_related_product_id" to table: "products"
CREATE INDEX "idx_products_related_product_id" ON "public"."products" ("related_product_id");
-- create index "idx_products_status" to table: "products"
CREATE INDEX "idx_products_status" ON "public"."products" ("status");
-- create "variants" table
CREATE TABLE "public"."variants" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "sku" character varying(100) NOT NULL,
  "barcode" character varying(100) NULL,
  "image_url" character varying(255) NULL,
  "price" numeric(15,2) NOT NULL,
  "compare_at_price" numeric(15,2) NULL,
  "cost_price" numeric(15,2) NULL,
  "stock_quantity" bigint NULL DEFAULT 0,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_variants" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_variants_stock_quantity" CHECK (stock_quantity >= 0)
);
-- create index "idx_variants_product_id" to table: "variants"
CREATE INDEX "idx_variants_product_id" ON "public"."variants" ("product_id");
-- create index "idx_variants_sku" to table: "variants"
CREATE UNIQUE INDEX "idx_variants_sku" ON "public"."variants" ("sku");
-- create "coupon_items" table
CREATE TABLE "public"."coupon_items" (
  "coupon_id" bigint NOT NULL,
  "variant_id" bigint NOT NULL,
  PRIMARY KEY ("coupon_id", "variant_id"),
  CONSTRAINT "fk_coupon_items_coupon" FOREIGN KEY ("coupon_id") REFERENCES "public"."coupons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_coupon_items_variant" FOREIGN KEY ("variant_id") REFERENCES "public"."variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "orders" table
CREATE TABLE "public"."orders" (
  "id" bigserial NOT NULL,
  "order_number" character varying(50) NOT NULL,
  "customer_email" character varying(255) NULL,
  "customer_phone" character varying(20) NULL,
  "total_price" numeric(15,2) NULL DEFAULT 0,
  "subtotal_price" numeric(15,2) NULL DEFAULT 0,
  "total_discounts" numeric(15,2) NULL DEFAULT 0,
  "total_tax" numeric(15,2) NULL DEFAULT 0,
  "shipping_cost" numeric(15,2) NULL DEFAULT 0,
  "currency" character varying(3) NULL DEFAULT 'USD',
  "exchange_rate" numeric(15,6) NULL DEFAULT 1,
  "financial_status" character varying(30) NULL,
  "fulfillment_status" character varying(30) NULL,
  "shipping_address" jsonb NULL,
  "billing_address" jsonb NULL,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_orders_deleted_at" to table: "orders"
CREATE INDEX "idx_orders_deleted_at" ON "public"."orders" ("deleted_at");
-- create index "idx_orders_number" to table: "orders"
CREATE UNIQUE INDEX "idx_orders_number" ON "public"."orders" ("order_number");
-- create "coupon_usages" table
CREATE TABLE "public"."coupon_usages" (
  "id" bigserial NOT NULL,
  "order_id" bigint NOT NULL,
  "coupon_id" bigint NOT NULL,
  "customer_email" character varying(255) NOT NULL,
  "discount_amount" numeric(15,2) NOT NULL,
  "used_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_coupon_usages_coupon" FOREIGN KEY ("coupon_id") REFERENCES "public"."coupons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_orders_coupon_usages" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_coupon_usages_deleted_at" to table: "coupon_usages"
CREATE INDEX "idx_coupon_usages_deleted_at" ON "public"."coupon_usages" ("deleted_at");
-- create "frequently_boughts" table
CREATE TABLE "public"."frequently_boughts" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "bought_with_product_id" bigint NOT NULL,
  "sort_order" bigint NULL DEFAULT 0,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_frequently_boughts_bought_with_product" FOREIGN KEY ("bought_with_product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_products_frequently_bought" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_product_bought_with" to table: "frequently_boughts"
CREATE UNIQUE INDEX "idx_product_bought_with" ON "public"."frequently_boughts" ("product_id", "bought_with_product_id");
-- create "options" table
CREATE TABLE "public"."options" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "name" character varying(100) NOT NULL,
  "code" character varying(100) NOT NULL,
  "position" bigint NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_options" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_options_product_id" to table: "options"
CREATE INDEX "idx_options_product_id" ON "public"."options" ("product_id");
-- create "option_values" table
CREATE TABLE "public"."option_values" (
  "id" bigserial NOT NULL,
  "option_id" bigint NOT NULL,
  "value" character varying(100) NOT NULL,
  "color_code" character varying(20) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_options_option_values" FOREIGN KEY ("option_id") REFERENCES "public"."options" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_option_values_option_id" to table: "option_values"
CREATE INDEX "idx_option_values_option_id" ON "public"."option_values" ("option_id");
-- create "order_items" table
CREATE TABLE "public"."order_items" (
  "id" bigserial NOT NULL,
  "order_id" bigint NOT NULL,
  "variant_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "product_name" character varying(255) NOT NULL,
  "variant_name" character varying(255) NULL,
  "sku" character varying(100) NULL,
  "quantity" bigint NOT NULL,
  "price" numeric(15,2) NOT NULL,
  "total" numeric(15,2) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_order_items_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_order_items_variant" FOREIGN KEY ("variant_id") REFERENCES "public"."variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_orders_order_items" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_order_items_quantity" CHECK (quantity > 0)
);
-- create index "idx_order_items_deleted_at" to table: "order_items"
CREATE INDEX "idx_order_items_deleted_at" ON "public"."order_items" ("deleted_at");
-- create index "idx_order_items_order_id" to table: "order_items"
CREATE INDEX "idx_order_items_order_id" ON "public"."order_items" ("order_id");
-- create index "idx_order_items_product_id" to table: "order_items"
CREATE INDEX "idx_order_items_product_id" ON "public"."order_items" ("product_id");
-- create index "idx_order_items_variant_id" to table: "order_items"
CREATE INDEX "idx_order_items_variant_id" ON "public"."order_items" ("variant_id");
-- create "policies" table
CREATE TABLE "public"."policies" (
  "id" bigserial NOT NULL,
  "country_code" character(2) NOT NULL,
  "type" character varying(50) NOT NULL,
  "title" character varying(255) NOT NULL,
  "content" text NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_policies_country" FOREIGN KEY ("country_code") REFERENCES "public"."countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_policy_country_type" to table: "policies"
CREATE INDEX "idx_policy_country_type" ON "public"."policies" ("country_code", "type");
-- create "product_categories" table
CREATE TABLE "public"."product_categories" (
  "category_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  PRIMARY KEY ("category_id", "product_id"),
  CONSTRAINT "fk_product_categories_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_product_categories_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "product_faqs" table
CREATE TABLE "public"."product_faqs" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "question" character varying(500) NOT NULL,
  "answer" text NOT NULL,
  "sort_order" bigint NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_fa_qs" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_product_faqs_product_id" to table: "product_faqs"
CREATE INDEX "idx_product_faqs_product_id" ON "public"."product_faqs" ("product_id");
-- create "product_images" table
CREATE TABLE "public"."product_images" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "image_url" character varying(500) NOT NULL,
  "video_url" character varying(500) NOT NULL,
  "media_type" character varying(50) NULL DEFAULT 'gallery',
  "highlight" boolean NULL DEFAULT false,
  "alt_text" character varying(255) NULL,
  "position" bigint NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_images" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create "suppliers" table
CREATE TABLE "public"."suppliers" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "platform" character varying(50) NOT NULL,
  "homepage_url" character varying(500) NULL,
  "contact_info" character varying(255) NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_suppliers_deleted_at" to table: "suppliers"
CREATE INDEX "idx_suppliers_deleted_at" ON "public"."suppliers" ("deleted_at");
-- create "product_mappings" table
CREATE TABLE "public"."product_mappings" (
  "id" bigserial NOT NULL,
  "local_variant_id" bigint NOT NULL,
  "supplier_id" bigint NOT NULL,
  "source_product_id" character varying(100) NOT NULL,
  "source_variant_id" character varying(100) NULL,
  "source_url" character varying(500) NOT NULL,
  "cost_price_cny" numeric(15,2) NULL DEFAULT 0,
  "cost_price_usd" numeric(15,2) NULL DEFAULT 0,
  "auto_sync_stock" boolean NULL DEFAULT true,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_product_mappings_local_variant" FOREIGN KEY ("local_variant_id") REFERENCES "public"."variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_product_mappings_supplier" FOREIGN KEY ("supplier_id") REFERENCES "public"."suppliers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_product_mappings_deleted_at" to table: "product_mappings"
CREATE INDEX "idx_product_mappings_deleted_at" ON "public"."product_mappings" ("deleted_at");
-- create index "uniq_mapping_local_var" to table: "product_mappings"
CREATE UNIQUE INDEX "uniq_mapping_local_var" ON "public"."product_mappings" ("local_variant_id");
-- create "product_price_tiers" table
CREATE TABLE "public"."product_price_tiers" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "name" character varying(255) NOT NULL,
  "qty" bigint NOT NULL DEFAULT 1,
  "savings_text" character varying(100) NULL,
  "price" numeric(15,2) NOT NULL,
  "base_price" numeric(15,2) NULL,
  "tag" character varying(50) NULL,
  "tag_class" character varying(100) NULL,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_price_tiers" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_product_price_tiers_deleted_at" to table: "product_price_tiers"
CREATE INDEX "idx_product_price_tiers_deleted_at" ON "public"."product_price_tiers" ("deleted_at");
-- create index "idx_product_price_tiers_product_id" to table: "product_price_tiers"
CREATE INDEX "idx_product_price_tiers_product_id" ON "public"."product_price_tiers" ("product_id");
-- create "product_reviews" table
CREATE TABLE "public"."product_reviews" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "author_name" character varying(100) NOT NULL,
  "author_email" character varying(255) NULL,
  "author_avatar" character varying(500) NULL,
  "rating" bigint NOT NULL,
  "content" text NULL,
  "media" jsonb NULL,
  "is_verified" boolean NULL DEFAULT false,
  "status" character varying(20) NULL DEFAULT 'pending',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_reviews" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_product_reviews_rating" CHECK ((rating >= 1) AND (rating <= 5))
);
-- create index "idx_product_reviews_author_email" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_author_email" ON "public"."product_reviews" ("author_email");
-- create index "idx_product_reviews_deleted_at" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_deleted_at" ON "public"."product_reviews" ("deleted_at");
-- create index "idx_product_reviews_product_id" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_product_id" ON "public"."product_reviews" ("product_id");
-- create index "idx_product_reviews_rating" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_rating" ON "public"."product_reviews" ("rating");
-- create index "idx_product_reviews_status" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_status" ON "public"."product_reviews" ("status");
-- create "purchase_orders" table
CREATE TABLE "public"."purchase_orders" (
  "id" bigserial NOT NULL,
  "supplier_id" bigint NOT NULL,
  "platform_order_id" character varying(100) NULL,
  "total_cost" numeric(15,2) NULL,
  "currency" character varying(3) NULL,
  "status" character varying(30) NULL,
  "local_tracking_number" character varying(100) NULL,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_purchase_orders_supplier" FOREIGN KEY ("supplier_id") REFERENCES "public"."suppliers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_purchase_orders_deleted_at" to table: "purchase_orders"
CREATE INDEX "idx_purchase_orders_deleted_at" ON "public"."purchase_orders" ("deleted_at");
-- create "purchase_order_items" table
CREATE TABLE "public"."purchase_order_items" (
  "id" bigserial NOT NULL,
  "purchase_order_id" bigint NOT NULL,
  "order_id" bigint NOT NULL,
  "variant_id" bigint NOT NULL,
  "quantity" bigint NOT NULL,
  "cost_per_item" numeric(15,2) NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_purchase_order_items_order" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_purchase_order_items_variant" FOREIGN KEY ("variant_id") REFERENCES "public"."variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_purchase_orders_purchase_order_items" FOREIGN KEY ("purchase_order_id") REFERENCES "public"."purchase_orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_purchase_order_items_deleted_at" to table: "purchase_order_items"
CREATE INDEX "idx_purchase_order_items_deleted_at" ON "public"."purchase_order_items" ("deleted_at");
-- create "shipments" table
CREATE TABLE "public"."shipments" (
  "id" bigserial NOT NULL,
  "order_id" bigint NOT NULL,
  "purchase_order_id" bigint NULL,
  "tracking_number" character varying(100) NULL,
  "carrier_code" character varying(50) NULL,
  "tracking_url" character varying(500) NULL,
  "status" character varying(30) NULL,
  "shipped_at" timestamptz NULL,
  "estimated_delivery_date" date NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_shipments" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_purchase_orders_shipments" FOREIGN KEY ("purchase_order_id") REFERENCES "public"."purchase_orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_shipments_deleted_at" to table: "shipments"
CREATE INDEX "idx_shipments_deleted_at" ON "public"."shipments" ("deleted_at");
-- create "sliders" table
CREATE TABLE "public"."sliders" (
  "id" bigserial NOT NULL,
  "title" character varying(255) NOT NULL,
  "image_url" character varying(500) NOT NULL,
  "sub_text" text NULL,
  "description" text NULL,
  "position" bigint NULL DEFAULT 0,
  "country_code" character(2) NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_sliders_country" FOREIGN KEY ("country_code") REFERENCES "public"."countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_sliders_country_code" to table: "sliders"
CREATE INDEX "idx_sliders_country_code" ON "public"."sliders" ("country_code");
-- create index "idx_sliders_deleted_at" to table: "sliders"
CREATE INDEX "idx_sliders_deleted_at" ON "public"."sliders" ("deleted_at");
-- create "transactions" table
CREATE TABLE "public"."transactions" (
  "id" bigserial NOT NULL,
  "order_id" bigint NOT NULL,
  "gateway" character varying(20) NOT NULL,
  "kind" character varying(20) NULL DEFAULT 'sale',
  "payment_method" character varying(50) NOT NULL,
  "transaction_reference" character varying(255) NULL,
  "amount" numeric(15,2) NOT NULL,
  "currency" character varying(3) NOT NULL,
  "status" character varying(20) NULL DEFAULT 'pending',
  "raw_response" jsonb NULL,
  "error_message" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_transactions" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_transactions_deleted_at" to table: "transactions"
CREATE INDEX "idx_transactions_deleted_at" ON "public"."transactions" ("deleted_at");
-- create "variant_value_map" table
CREATE TABLE "public"."variant_value_map" (
  "variant_id" bigint NOT NULL,
  "option_value_id" bigint NOT NULL,
  PRIMARY KEY ("variant_id", "option_value_id"),
  CONSTRAINT "fk_variant_value_map_option_value" FOREIGN KEY ("option_value_id") REFERENCES "public"."option_values" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_variant_value_map_variant" FOREIGN KEY ("variant_id") REFERENCES "public"."variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
