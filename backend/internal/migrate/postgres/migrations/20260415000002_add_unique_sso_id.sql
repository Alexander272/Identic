-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users ADD CONSTRAINT users_sso_id_key UNIQUE (sso_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.users DROP CONSTRAINT IF EXISTS users_sso_id_key;
-- +goose StatementEnd
