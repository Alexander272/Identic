-- +goose Up
-- +goose StatementBegin
-- 1. GIN индекс для нечеткого поиска по триграммам
CREATE INDEX IF NOT EXISTS idx_positions_name_trgm ON positions USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_positions_search_trgm ON positions USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_positions_search_exact ON positions USING btree (search);

-- 2. Обычный индекс для быстрого джойна по order_id
CREATE INDEX idx_positions_order_id 
ON positions (order_id);

-- 3. Индекс по году для быстрой группировки и фильтрации
CREATE INDEX idx_orders_year ON orders(year DESC);

-- 4. Составной индекс (год + id), если часто выбираем конкретные года
CREATE INDEX idx_orders_year_id ON orders(year, id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_positions_order_id;
DROP INDEX IF EXISTS idx_positions_search_trgm;
DROP INDEX IF EXISTS idx_positions_name_trgm;
DROP INDEX IF EXISTS idx_positions_search_exact;
DROP INDEX IF EXISTS idx_orders_year;
DROP INDEX IF EXISTS idx_orders_year_id;
-- +goose StatementEnd
