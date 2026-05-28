package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
)

const tablePrices = "prices"

type PricesRepo struct {
	db *pgxpool.Pool
}

func NewPricesRepo(db *pgxpool.Pool) *PricesRepo {
	return &PricesRepo{db: db}
}

type Prices interface {
	Search(ctx context.Context, query string, codes []string, page, perPage int) ([]*models.Price, int, error)
	SearchAll(ctx context.Context, query string, codes []string) ([]*models.Price, error)
	CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.Price) error
	UpsertSeveral(ctx context.Context, tx postgres.Tx, dto []*models.Price) error
	DeleteSeveral(ctx context.Context, tx postgres.Tx, codes []string) error
}

func (r *PricesRepo) getExec(tx postgres.Tx) postgres.QueryExecutor {
	if tx != nil {
		return tx.TX()
	}
	return r.db
}

func (r *PricesRepo) Search(ctx context.Context, query string, codes []string, page, perPage int) ([]*models.Price, int, error) {
	columns := "id, code, current_name, new_name, price, template, note, technique, under_drawing"

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ($1 = '' OR search_text ILIKE '%%' || $1 || '%%')
		AND (cardinality($2::text[]) = 0 OR code = ANY($2))`,
		tablePrices,
	)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, query, codes).Scan(&total); err != nil {
		return nil, 0, postgres.MapError(fmt.Errorf("failed to count: %w", err))
	}

	dataQuery := fmt.Sprintf(`
		SELECT %s FROM %s
		WHERE ($1 = '' OR search_text ILIKE '%%' || $1 || '%%')
		  AND (cardinality($2::text[]) = 0 OR code = ANY($2))
		ORDER BY
			CASE WHEN $1 = '' THEN 0 ELSE Prices($1 IN search_text) END,
			array_position($2::text[], code)
		LIMIT $3 OFFSET $4`, columns, tablePrices)

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx, dataQuery, query, codes, perPage, offset)
	if err != nil {
		return nil, 0, postgres.MapError(fmt.Errorf("failed to execute query: %w", err))
	}
	defer rows.Close()

	positions, err := scanPositions(rows)
	if err != nil {
		return nil, 0, err
	}
	return positions, total, nil
}

func (r *PricesRepo) SearchAll(ctx context.Context, query string, codes []string) ([]*models.Price, error) {
	columns := "id, code, current_name, new_name, price, template, note, technique, under_drawing"

	queryStr := fmt.Sprintf(`
		SELECT %s FROM %s
		WHERE ($1 = '' OR search_text ILIKE '%%' || $1 || '%%')
		  AND (cardinality($2::text[]) = 0 OR code = ANY($2))
		ORDER BY
			CASE WHEN $1 = '' THEN 0 ELSE Prices($1 IN search_text) END,
			array_position($2::text[], code)`, columns, tablePrices)

	rows, err := r.db.Query(ctx, queryStr, query, codes)
	if err != nil {
		return nil, postgres.MapError(fmt.Errorf("failed to execute query: %w", err))
	}
	defer rows.Close()

	return scanPositions(rows)
}

func (r *PricesRepo) CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.Price) error {
	if len(dto) == 0 {
		return nil
	}

	rows := make([][]interface{}, len(dto))
	for i, p := range dto {
		searchText := buildSearchText(p)
		rows[i] = []interface{}{p.Code, p.CurrentName, p.NewName, p.Price, p.Template, p.Note, p.Technique, p.UnderDrawing, searchText}
	}

	columns := []string{"code", "current_name", "new_name", "price", "template", "note", "technique", "under_drawing", "search_text"}
	_, err := r.getExec(tx).CopyFrom(ctx, pgx.Identifier{tablePrices}, columns, pgx.CopyFromRows(rows))
	if err != nil {
		return postgres.MapError(fmt.Errorf("failed to execute query: %w", err))
	}
	return nil
}

func (r *PricesRepo) UpsertSeveral(ctx context.Context, tx postgres.Tx, dto []*models.Price) error {
	if len(dto) == 0 {
		return nil
	}

	sql := fmt.Sprintf(`
		INSERT INTO %s (code, current_name, new_name, price, template, note, technique, under_drawing, search_text)
		SELECT * FROM UNNEST($1::text[], $2::text[], $3::text[], $4::float8[], $5::text[], $6::text[], $7::text[], $8::text[], $9::text[])
		ON CONFLICT (code) DO UPDATE SET
			current_name = EXCLUDED.current_name,
			new_name = EXCLUDED.new_name,
			price = EXCLUDED.price,
			template = EXCLUDED.template,
			note = EXCLUDED.note,
			technique = EXCLUDED.technique,
			under_drawing = EXCLUDED.under_drawing,
			search_text = EXCLUDED.search_text
	`, tablePrices)

	codes := make([]string, len(dto))
	currentNames := make([]string, len(dto))
	newNames := make([]string, len(dto))
	prices := make([]float64, len(dto))
	templates := make([]string, len(dto))
	notes := make([]string, len(dto))
	techniques := make([]string, len(dto))
	underDrawings := make([]string, len(dto))
	searchTexts := make([]string, len(dto))

	for i, p := range dto {
		codes[i] = p.Code
		currentNames[i] = p.CurrentName
		newNames[i] = p.NewName
		prices[i] = p.Price
		templates[i] = p.Template
		notes[i] = p.Note
		techniques[i] = p.Technique
		underDrawings[i] = p.UnderDrawing
		searchTexts[i] = buildSearchText(p)
	}

	_, err := r.getExec(tx).Exec(ctx, sql,
		codes, currentNames, newNames, prices, templates, notes, techniques, underDrawings, searchTexts,
	)
	if err != nil {
		return postgres.MapError(fmt.Errorf("failed to execute query: %w", err))
	}
	return nil
}

func (r *PricesRepo) DeleteSeveral(ctx context.Context, tx postgres.Tx, codes []string) error {
	if len(codes) == 0 {
		return nil
	}

	sql := fmt.Sprintf(`DELETE FROM %s WHERE code = ANY($1::text[])`, tablePrices)
	_, err := r.getExec(tx).Exec(ctx, sql, codes)
	if err != nil {
		return postgres.MapError(fmt.Errorf("failed to execute query: %w", err))
	}
	return nil
}

func scanPositions(rows pgx.Rows) ([]*models.Price, error) {
	var data []*models.Price
	for rows.Next() {
		var p models.Price
		if err := rows.Scan(&p.ID, &p.Code, &p.CurrentName, &p.NewName, &p.Price, &p.Template, &p.Note, &p.Technique, &p.UnderDrawing); err != nil {
			return nil, postgres.MapError(fmt.Errorf("scan row error: %w", err))
		}
		data = append(data, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, postgres.MapError(fmt.Errorf("rows iteration error: %w", err))
	}
	return data, nil
}

func buildSearchText(p *models.Price) string {
	parts := []string{p.CurrentName, p.NewName, strconv.FormatFloat(p.Price, 'f', -1, 64), p.Template}
	text := strings.Join(parts, " ")
	return normalizeSearch(text)
}

func normalizeSearch(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "х", "x")
	return s
}
