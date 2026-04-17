-- +goose Up
-- +goose StatementBegin
ALTER TABLE activity_logs ADD COLUMN IF NOT EXISTS entity TEXT COLLATE pg_catalog."default";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE activity_logs DROP COLUMN IF EXISTS entity;
-- +goose StatementEnd
