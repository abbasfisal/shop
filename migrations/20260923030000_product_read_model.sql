-- +goose Up

-- flattened product document (replaces the MongoDB products collection)
ALTER TABLE products ADD COLUMN read_model JSONB DEFAULT '{}';

-- product recommendations (replaces the MongoDB recommendations collection)
CREATE TABLE product_recommendations (
    id                     BIGSERIAL PRIMARY KEY,
    product_id             BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    recommended_product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at             TIMESTAMPTZ DEFAULT now(),
    updated_at             TIMESTAMPTZ DEFAULT now(),
    UNIQUE (product_id, recommended_product_id)
);
CREATE INDEX idx_product_recommendations_product_id ON product_recommendations (product_id);

-- +goose Down
DROP TABLE IF EXISTS product_recommendations;
ALTER TABLE products DROP COLUMN IF EXISTS read_model;
