-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.user_logins DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE public.user_logins ALTER COLUMN user_id TYPE TEXT;
ALTER TABLE public.user_logins ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(sso_id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.user_logins DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE public.user_logins ALTER COLUMN user_id TYPE UUID USING user_id::UUID;
ALTER TABLE public.user_logins ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd
