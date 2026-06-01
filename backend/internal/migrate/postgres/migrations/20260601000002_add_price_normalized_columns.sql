-- +goose Up
-- +goose StatementBegin
ALTER TABLE prices ADD COLUMN IF NOT EXISTS current_name_norm TEXT NOT NULL DEFAULT '';
ALTER TABLE prices ADD COLUMN IF NOT EXISTS new_name_norm TEXT NOT NULL DEFAULT '';
ALTER TABLE prices ADD COLUMN IF NOT EXISTS template_norm TEXT NOT NULL DEFAULT '';

UPDATE prices SET 
    current_name_norm = REPLACE(LOWER(current_name), 'х', 'x'),
    new_name_norm = REPLACE(LOWER(new_name), 'х', 'x'),
    template_norm = REPLACE(LOWER(template), 'х', 'x');

CREATE INDEX IF NOT EXISTS idx_prices_cn_norm ON prices USING GIN(current_name_norm gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_prices_nn_norm ON prices USING GIN(new_name_norm gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_prices_tpl_norm ON prices USING GIN(template_norm gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_prices_cn_norm;
DROP INDEX IF EXISTS idx_prices_nn_norm;
DROP INDEX IF EXISTS idx_prices_tpl_norm;

ALTER TABLE prices DROP COLUMN IF EXISTS current_name_norm;
ALTER TABLE prices DROP COLUMN IF EXISTS new_name_norm;
ALTER TABLE prices DROP COLUMN IF EXISTS template_norm;
-- +goose StatementEnd
