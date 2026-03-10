-- +goose Up
-- +goose StatementBegin
-- 1. GIN индекс для нечеткого поиска по триграммам
CREATE INDEX IF NOT EXISTS idx_positions_name_trgm ON positions USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_positions_search_trgm ON positions USING gin (name gin_trgm_ops);

-- 2. Обычный индекс для быстрого джойна по order_id
CREATE INDEX idx_positions_order_id 
ON positions (order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_positions_search_trgm;
DROP INDEX IF EXISTS idx_positions_name_trgm;
-- +goose StatementEnd
