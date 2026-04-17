-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS search_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    search_id TEXT COLLATE pg_catalog."default" NOT NULL,
    actor_id UUID NOT NULL,
    actor_name TEXT COLLATE pg_catalog."default" NOT NULL,
    search_type TEXT NOT NULL,                      -- 'exact' или 'fuzzy'
    query JSONB NOT NULL,                           -- Сохранённый запрос
    duration_ms BIGINT NOT NULL,                    -- Длительность в миллисекундах
    results_count INTEGER NOT NULL,                 -- Количество найденных результатов
    items_count INTEGER NOT NULL,                   -- Количество позиций в запросе
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.search_logs
    OWNER TO postgres;

CREATE INDEX IF NOT EXISTS idx_search_logs_actor ON search_logs (actor_id);
CREATE INDEX IF NOT EXISTS idx_search_logs_created ON search_logs (created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.search_logs;
-- +goose StatementEnd
