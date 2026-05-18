-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.user_logins ALTER COLUMN user_id DROP NOT NULL;
ALTER TABLE public.user_logins DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE public.user_logins ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(sso_id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.user_logins DROP CONSTRAINT IF EXISTS fk_user;
DELETE FROM public.user_logins WHERE user_id IS NULL;
ALTER TABLE public.user_logins ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE public.user_logins ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(sso_id) ON DELETE CASCADE;
-- +goose StatementEnd
