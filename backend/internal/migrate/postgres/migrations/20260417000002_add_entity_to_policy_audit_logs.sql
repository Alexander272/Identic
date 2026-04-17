-- +goose Up
-- +goose StatementBegin
ALTER TABLE policy_audit_logs ADD COLUMN IF NOT EXISTS entity TEXT COLLATE pg_catalog."default";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE policy_audit_logs DROP COLUMN IF EXISTS entity;
-- +goose StatementEnd
