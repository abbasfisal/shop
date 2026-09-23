-- +goose Up

-- ============ users ============
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    first_name    VARCHAR(50) NOT NULL,
    last_name     VARCHAR(50) NOT NULL,
    phone_number  VARCHAR(11) NOT NULL,
    password      TEXT,
    type          VARCHAR(10) NOT NULL DEFAULT 'client' CHECK (type IN ('admin','client')),
    created_at    TIMESTAMPTZ DEFAULT now(),
    updated_at    TIMESTAMPTZ DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- ============ customers ============
CREATE TABLE customers (
    id          BIGSERIAL PRIMARY KEY,
    mobile      TEXT,
    first_name  TEXT,
    last_name   TEXT,
    active      BOOLEAN,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX "unique" ON customers (mobile);
CREATE INDEX idx_customers_deleted_at ON customers (deleted_at);

-- ============ addresses ============
CREATE TABLE addresses (
    id                  BIGSERIAL PRIMARY KEY,
    customer_id         BIGINT NOT NULL REFERENCES customers(id),
    receiver_name       TEXT,
    receiver_mobile     TEXT,
    receiver_address    TEXT,
    receiver_postal_code TEXT,
    created_at          TIMESTAMPTZ DEFAULT now(),
    updated_at          TIMESTAMPTZ DEFAULT now(),
    deleted_at          TIMESTAMPTZ
);
CREATE INDEX idx_addresses_deleted_at ON addresses (deleted_at);

-- ============ brands ============
CREATE TABLE brands (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(150) NOT NULL,
    slug        VARCHAR(150) NOT NULL UNIQUE,
    image       TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_brands_deleted_at ON brands (deleted_at);

-- ============ categories (self-referencing tree) ============
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    priority    INTEGER,
    title       VARCHAR(150) NOT NULL,
    slug        VARCHAR(150) NOT NULL UNIQUE,
    parent_id   BIGINT REFERENCES categories(id),
    image       TEXT,
    status      BOOLEAN,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_categories_parent_id ON categories (parent_id);
CREATE INDEX idx_categories_deleted_at ON categories (deleted_at);

-- ============ attributes ============
CREATE TABLE attributes (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_attributes_deleted_at ON attributes (deleted_at);

-- ============ attribute_values ============
CREATE TABLE attribute_values (
    id              BIGSERIAL PRIMARY KEY,
    attribute_id    BIGINT NOT NULL REFERENCES attributes(id),
    attribute_title TEXT,
    value           TEXT,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_attribute_values_attribute_id ON attribute_values (attribute_id);
CREATE INDEX idx_attribute_values_deleted_at ON attribute_values (deleted_at);

-- ============ products ============
CREATE TABLE products (
    id              BIGSERIAL PRIMARY KEY,
    category_id     BIGINT REFERENCES categories(id),
    brand_id        BIGINT REFERENCES brands(id),
    title           TEXT,
    slug            TEXT UNIQUE,
    sku             TEXT UNIQUE,
    status          BOOLEAN,
    original_price  BIGINT,
    sale_price      BIGINT,
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_products_category_id ON products (category_id);
CREATE INDEX idx_products_brand_id ON products (brand_id);
CREATE INDEX idx_products_deleted_at ON products (deleted_at);

-- ============ product_images ============
CREATE TABLE product_images (
    id          BIGSERIAL PRIMARY KEY,
    product_id  BIGINT NOT NULL REFERENCES products(id),
    path        TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_product_images_product_id ON product_images (product_id);
CREATE INDEX idx_product_images_deleted_at ON product_images (deleted_at);

-- ============ features ============
CREATE TABLE features (
    id          BIGSERIAL PRIMARY KEY,
    product_id  BIGINT NOT NULL REFERENCES products(id),
    title       TEXT,
    value       TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_features_product_id ON features (product_id);
CREATE INDEX idx_features_deleted_at ON features (deleted_at);

-- ============ product_attributes ============
CREATE TABLE product_attributes (
    id                     BIGSERIAL PRIMARY KEY,
    product_id             BIGINT NOT NULL REFERENCES products(id),
    attribute_id           BIGINT NOT NULL REFERENCES attributes(id),
    attribute_title        TEXT,
    attribute_value_id     BIGINT NOT NULL REFERENCES attribute_values(id),
    attribute_value_title  TEXT,
    created_at             TIMESTAMPTZ DEFAULT now(),
    updated_at             TIMESTAMPTZ DEFAULT now(),
    deleted_at             TIMESTAMPTZ
);
CREATE INDEX idx_product_attributes_product_id ON product_attributes (product_id);
CREATE INDEX idx_product_attributes_attribute_id ON product_attributes (attribute_id);
CREATE INDEX idx_product_attributes_deleted_at ON product_attributes (deleted_at);

-- ============ product_inventories ============
CREATE TABLE product_inventories (
    id              BIGSERIAL PRIMARY KEY,
    product_id      BIGINT NOT NULL REFERENCES products(id),
    quantity        BIGINT,
    reserved_stock  BIGINT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_product_inventories_product_id ON product_inventories (product_id);
CREATE INDEX idx_product_inventories_deleted_at ON product_inventories (deleted_at);

-- ============ product_inventory_attributes ============
CREATE TABLE product_inventory_attributes (
    id                    BIGSERIAL PRIMARY KEY,
    product_id            BIGINT NOT NULL REFERENCES products(id),
    product_inventory_id  BIGINT NOT NULL REFERENCES product_inventories(id),
    product_attribute_id  BIGINT NOT NULL REFERENCES product_attributes(id),
    created_at            TIMESTAMPTZ DEFAULT now(),
    updated_at            TIMESTAMPTZ DEFAULT now(),
    deleted_at            TIMESTAMPTZ
);
CREATE INDEX idx_pia_product_id ON product_inventory_attributes (product_id);
CREATE INDEX idx_pia_product_inventory_id ON product_inventory_attributes (product_inventory_id);
CREATE INDEX idx_pia_product_attribute_id ON product_inventory_attributes (product_attribute_id);
CREATE INDEX idx_pia_deleted_at ON product_inventory_attributes (deleted_at);

-- ============ carts ============
CREATE TABLE carts (
    id           BIGSERIAL PRIMARY KEY,
    customer_id  BIGINT NOT NULL REFERENCES customers(id),
    status       SMALLINT,
    created_at   TIMESTAMPTZ DEFAULT now(),
    updated_at   TIMESTAMPTZ DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX idx_carts_customer_id ON carts (customer_id);
CREATE INDEX idx_carts_deleted_at ON carts (deleted_at);

-- ============ cart_items ============
CREATE TABLE cart_items (
    id              BIGSERIAL PRIMARY KEY,
    customer_id     BIGINT NOT NULL REFERENCES customers(id),
    cart_id         BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id      BIGINT NOT NULL REFERENCES products(id),
    inventory_id    BIGINT,
    quantity        SMALLINT,
    original_price  BIGINT,
    sale_price      BIGINT,
    product_sku     TEXT,
    product_title   TEXT,
    product_image   TEXT,
    product_slug    TEXT,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_cart_items_customer_id ON cart_items (customer_id);
CREATE INDEX idx_cart_items_cart_id ON cart_items (cart_id);
CREATE INDEX idx_cart_items_product_id ON cart_items (product_id);
CREATE INDEX idx_cart_items_deleted_at ON cart_items (deleted_at);

-- ============ orders ============
CREATE TABLE orders (
    id                    BIGSERIAL PRIMARY KEY,
    customer_id           BIGINT NOT NULL REFERENCES customers(id),
    order_number          VARCHAR(10) UNIQUE,
    payment_status        INTEGER,
    total_original_price  BIGINT,
    total_sale_price      BIGINT,
    discount              BIGINT,
    order_status          INTEGER,
    address               TEXT,
    note                  TEXT,
    created_at            TIMESTAMPTZ DEFAULT now(),
    updated_at            TIMESTAMPTZ DEFAULT now(),
    deleted_at            TIMESTAMPTZ
);
CREATE INDEX idx_orders_customer_id ON orders (customer_id);
CREATE INDEX idx_orders_deleted_at ON orders (deleted_at);

-- ============ order_items ============
CREATE TABLE order_items (
    id                    BIGSERIAL PRIMARY KEY,
    customer_id           BIGINT NOT NULL REFERENCES customers(id),
    order_id              BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id            BIGINT NOT NULL REFERENCES products(id),
    inventory_id          BIGINT,
    quantity              BIGINT,
    original_price        BIGINT,
    sale_price            BIGINT,
    total_original_price  BIGINT,
    total_sale_price      BIGINT,
    created_at            TIMESTAMPTZ DEFAULT now(),
    updated_at            TIMESTAMPTZ DEFAULT now(),
    deleted_at            TIMESTAMPTZ
);
CREATE INDEX idx_order_items_order_id ON order_items (order_id);
CREATE INDEX idx_order_items_product_id ON order_items (product_id);
CREATE INDEX idx_order_items_deleted_at ON order_items (deleted_at);

-- ============ payments ============
CREATE TABLE payments (
    id           BIGSERIAL PRIMARY KEY,
    customer_id  BIGINT NOT NULL REFERENCES customers(id),
    order_id     BIGINT NOT NULL REFERENCES orders(id),
    authority    TEXT UNIQUE,
    description  TEXT,
    payment_url  TEXT,
    status_code  INTEGER,
    amount       BIGINT,
    ref_id       TEXT,
    status       INTEGER,
    created_at   TIMESTAMPTZ DEFAULT now(),
    updated_at   TIMESTAMPTZ DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX idx_payments_order_id ON payments (order_id);
CREATE INDEX idx_payments_deleted_at ON payments (deleted_at);

-- ============ otps ============
CREATE TABLE otps (
    id          BIGSERIAL PRIMARY KEY,
    mobile      TEXT,
    code        TEXT,
    is_expired  BOOLEAN,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_otps_deleted_at ON otps (deleted_at);

-- ============ sessions ============
CREATE TABLE sessions (
    id           BIGSERIAL PRIMARY KEY,
    mobile       TEXT,
    customer_id  BIGINT NOT NULL REFERENCES customers(id),
    session_id   TEXT,
    is_active    BOOLEAN,
    expired_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT now(),
    updated_at   TIMESTAMPTZ DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX idx_sessions_customer_id ON sessions (customer_id);
CREATE INDEX idx_sessions_deleted_at ON sessions (deleted_at);

-- ============ banners ============
CREATE TABLE banners (
    id          BIGSERIAL PRIMARY KEY,
    type        INTEGER,
    link        TEXT,
    priority    INTEGER,
    status      BOOLEAN,
    image       TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_banners_deleted_at ON banners (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS otps;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS product_inventory_attributes;
DROP TABLE IF EXISTS product_inventories;
DROP TABLE IF EXISTS product_attributes;
DROP TABLE IF EXISTS features;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS attribute_values;
DROP TABLE IF EXISTS attributes;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS users;
