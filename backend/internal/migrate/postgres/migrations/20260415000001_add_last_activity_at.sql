-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.user_logins ADD COLUMN last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
CREATE INDEX idx_user_logins_last_activity ON user_logins (user_id, last_activity_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_logins_last_activity;
ALTER TABLE public.user_logins DROP COLUMN IF EXISTS last_activity_at;
-- +goose StatementEnd
