-- +goose Up
-- ============================================================
-- Laravel-style product status enum (draft | published | archived)
-- + product expires_at (Laravel products.expires_at)
-- ============================================================
ALTER TABLE products ALTER COLUMN status DROP DEFAULT;
ALTER TABLE products ALTER COLUMN status TYPE VARCHAR(16)
    USING (CASE WHEN status THEN 'published' ELSE 'draft' END);
ALTER TABLE products ALTER COLUMN status SET DEFAULT 'draft';
ALTER TABLE products ADD COLUMN expires_at DATE;

-- +goose Down
ALTER TABLE products DROP COLUMN IF EXISTS expires_at;
ALTER TABLE products ALTER COLUMN status DROP DEFAULT;
ALTER TABLE products ALTER COLUMN status TYPE BOOLEAN
    USING (status = 'published');
ALTER TABLE products ALTER COLUMN status SET DEFAULT FALSE;
