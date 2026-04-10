-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_role_permissions_role;
-- +goose StatementEnd
