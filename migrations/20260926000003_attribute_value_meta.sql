-- +goose Up
-- ============================================================
-- Dynamic attribute values:
--   attributes.input_type  (text | color) — decides which input the
--     admin value-form renders (plain text vs color picker)
--   attribute_values.meta  (JSONB) — generic extras; colors store
--     {"hex": "#rrggbb"} so the storefront can paint swatches.
-- Backfills the well-known color attribute of the demo dataset.
-- ============================================================

ALTER TABLE attributes ADD COLUMN input_type VARCHAR(16) NOT NULL DEFAULT 'text';
ALTER TABLE attributes ADD CONSTRAINT chk_attributes_input_type CHECK (input_type IN ('text', 'color'));

ALTER TABLE attribute_values ADD COLUMN meta JSONB NOT NULL DEFAULT '{}';
CREATE INDEX idx_attribute_values_meta ON attribute_values USING GIN (meta);

-- the color attribute (stable key: code) becomes a color picker
UPDATE attributes SET input_type = 'color' WHERE code = 'color';

-- hex for the seeded color values (attribute code + value = stable key)
UPDATE attribute_values SET meta = '{"hex": "#2563eb"}'
WHERE value = 'آبی'
  AND attribute_id IN (SELECT id FROM attributes WHERE code = 'color');

UPDATE attribute_values SET meta = '{"hex": "#dc2626"}'
WHERE value = 'قرمز'
  AND attribute_id IN (SELECT id FROM attributes WHERE code = 'color');

UPDATE attribute_values SET meta = '{"hex": "#7c3aed"}'
WHERE value = 'بنفش'
  AND attribute_id IN (SELECT id FROM attributes WHERE code = 'color');

-- +goose Down
UPDATE attribute_values SET meta = '{}' WHERE meta <> '{}';
UPDATE attributes SET input_type = 'text' WHERE input_type <> 'text';
DROP INDEX IF EXISTS idx_attribute_values_meta;
ALTER TABLE attribute_values DROP COLUMN IF EXISTS meta;
ALTER TABLE attributes DROP CONSTRAINT IF EXISTS chk_attributes_input_type;
ALTER TABLE attributes DROP COLUMN IF EXISTS input_type;
