-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.price_search_logs
    ADD COLUMN IF NOT EXISTS fields JSONB NOT NULL DEFAULT '[]';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.price_search_logs
    DROP COLUMN IF EXISTS fields;
-- +goose StatementEnd
