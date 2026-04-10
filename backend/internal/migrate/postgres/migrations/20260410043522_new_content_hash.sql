-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN IF NOT EXISTS content_hash TEXT;

CREATE INDEX IF NOT EXISTS idx_orders_customer_hash_created 
ON orders (customer, content_hash, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_customer_hash_created;
ALTER TABLE orders DROP COLUMN IF EXISTS content_hash;
-- +goose StatementEnd
