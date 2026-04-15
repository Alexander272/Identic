-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.user_logins (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    login_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB, -- Для доп. инфо: геопозиция, ID сессии и т.д.
    
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.user_logins
    OWNER to postgres;

CREATE INDEX idx_user_logins_user_id_at ON user_logins (user_id, login_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_logins_user_id_at;

DROP TABLE IF EXISTS user_logins;
-- +goose StatementEnd
