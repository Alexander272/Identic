-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS price_search_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    queries JSONB NOT NULL,
    codes JSONB NOT NULL DEFAULT '[]',
    actor_id UUID NOT NULL,
    actor_name TEXT NOT NULL DEFAULT '',
    results_count INTEGER NOT NULL DEFAULT 0,
    duration_ms BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.price_search_logs
    OWNER TO postgres;

CREATE INDEX IF NOT EXISTS idx_price_search_logs_actor ON price_search_logs (actor_id);
CREATE INDEX IF NOT EXISTS idx_price_search_logs_created ON price_search_logs (created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.price_search_logs;
-- +goose StatementEnd
