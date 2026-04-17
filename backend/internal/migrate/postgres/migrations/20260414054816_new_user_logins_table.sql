-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.user_logins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    login_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB
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
