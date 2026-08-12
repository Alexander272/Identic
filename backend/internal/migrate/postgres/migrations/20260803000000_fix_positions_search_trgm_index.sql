-- +goose Up
-- +goose StatementBegin
-- Старый индекс ошибочно построен на name (а не на search): из-за этого и точный
-- (LIKE '%...%'), и нечёткий (p.search % q) поиск идут полным сканом по positions.
DROP INDEX IF EXISTS idx_positions_search_trgm;
CREATE INDEX IF NOT EXISTS idx_positions_search_trgm ON positions USING gin (search gin_trgm_ops);
-- Индекс для поиска по примечаниям (точный поиск матчит и normalized_notes)
CREATE INDEX IF NOT EXISTS idx_positions_notes_trgm ON positions USING gin (normalized_notes gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_positions_notes_trgm;
DROP INDEX IF EXISTS idx_positions_search_trgm;
CREATE INDEX IF NOT EXISTS idx_positions_search_trgm ON positions USING gin (name gin_trgm_ops);
-- +goose StatementEnd
