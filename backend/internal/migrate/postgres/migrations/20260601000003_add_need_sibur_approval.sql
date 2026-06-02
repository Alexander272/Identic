-- +goose Up
-- +goose StatementBegin
ALTER TABLE prices ADD COLUMN need_sibur_approval TEXT NOT NULL DEFAULT '';
ALTER TABLE prices DROP COLUMN technique;
ALTER TABLE prices DROP COLUMN under_drawing;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE prices ADD COLUMN technique TEXT NOT NULL DEFAULT '';
ALTER TABLE prices ADD COLUMN under_drawing TEXT NOT NULL DEFAULT '';
ALTER TABLE prices DROP COLUMN need_sibur_approval;
-- +goose StatementEnd
