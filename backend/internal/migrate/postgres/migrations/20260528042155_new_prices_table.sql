-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

DROP TABLE IF EXISTS prices;

CREATE TABLE prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL,
    current_name TEXT NOT NULL DEFAULT '',
    new_name TEXT NOT NULL DEFAULT '',
    price NUMERIC(12,2) NOT NULL DEFAULT 0,
    template TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    technique TEXT NOT NULL DEFAULT '',
    under_drawing TEXT NOT NULL DEFAULT '',
    search_text TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_prices_code ON prices(code);
CREATE INDEX idx_prices_search ON prices USING GIN(search_text gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS prices;

DROP INDEX IF EXISTS idx_prices_code;
DROP INDEX IF EXISTS idx_prices_search;
-- +goose StatementEnd
