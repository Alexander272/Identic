-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action TEXT NOT NULL,                                   -- 'INSERT', 'UPDATE', 'DELETE'
    changed_by TEXT COLLATE pg_catalog."default" NOT NULL,  -- Кто изменил (UserID)
    changed_by_name TEXT COLLATE pg_catalog."default" NOT NULL,
    
    entity_type TEXT NOT NULL,          -- 'order' или 'order_item'
    entity_id UUID NOT NULL,            -- ID заказа или позиции
    parent_id UUID,                     -- ID заказа (для позиций, чтобы легко найти все изменения внутри заказа)
   
    old_values JSONB,                   -- Состояние ДО
    new_values JSONB,                   -- Состояние ПОСЛЕ
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.activity_logs
    OWNER to postgres;

CREATE INDEX IF NOT EXISTS idx_activity_order ON activity_logs (parent_id) WHERE parent_id IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.activity_logs;
-- +goose StatementEnd
