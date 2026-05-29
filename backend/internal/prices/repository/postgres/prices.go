package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
)

type PricesRepo struct {
	db *pgxpool.Pool
}

func NewPricesRepo(db *pgxpool.Pool) *PricesRepo {
	return &PricesRepo{db: db}
}

type Prices interface {
	Search(ctx context.Context, queries []string, codes []string, page, perPage int) ([]*models.Price, int, error)
	SearchAll(ctx context.Context, queries []string, codes []string) ([]*models.Price, error)
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

func (r *PricesRepo) Search(ctx context.Context, queries []string, codes []string, page, perPage int) ([]*models.Price, int, error) {
	columns := "id, code, current_name, new_name, price, template, note, technique, under_drawing"

	var conditions []string
	args := make([]any, 0)

	var firstQueryIdx int
	if len(queries) > 0 {
		firstQueryIdx = len(args) + 1
		queryConds := make([]string, len(queries))
		for i, q := range queries {
			queryConds[i] = fmt.Sprintf("search_text ILIKE '%%' || $%d || '%%'", i+1)
			args = append(args, q)
		}
		conditions = append(conditions, strings.Join(queryConds, " AND "))
	}

	var codesParamIdx int
	if len(codes) > 0 {
		codesParamIdx = len(args) + 1
		conditions = append(conditions, fmt.Sprintf("code::text = ANY($%d::text[])", codesParamIdx))
		args = append(args, codes)
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", Tables.Prices, where)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, postgres.MapError(fmt.Errorf("failed to count: %w", err))
	}

	var orderConds []string
	if len(queries) == 1 {
		orderConds = append(orderConds, fmt.Sprintf("similarity(search_text, $%d) DESC", firstQueryIdx))
	}
	if len(codes) > 0 {
		orderConds = append(orderConds, fmt.Sprintf("array_position($%d::text[], code::text)", codesParamIdx))
	}

	orderBy := ""
	if len(orderConds) > 0 {
		orderBy = "ORDER BY " + strings.Join(orderConds, ", ")
	}

	offset := (page - 1) * perPage
	perPageIdx := len(args) + 1
	offsetIdx := len(args) + 2
	dataArgs := make([]any, len(args)+2)
	copy(dataArgs, args)
	dataArgs[len(args)] = perPage
	dataArgs[len(args)+1] = offset

	dataQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s %s LIMIT $%d OFFSET $%d",
		columns, Tables.Prices, where, orderBy, perPageIdx, offsetIdx)

	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
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

func (r *PricesRepo) SearchAll(ctx context.Context, queries []string, codes []string) ([]*models.Price, error) {
	columns := "id, code, current_name, new_name, price, template, note, technique, under_drawing"

	var conditions []string
	args := make([]any, 0)

	var firstQueryIdx int
	if len(queries) > 0 {
		firstQueryIdx = len(args) + 1
		queryConds := make([]string, len(queries))
		for i, q := range queries {
			queryConds[i] = fmt.Sprintf("search_text ILIKE '%%' || $%d || '%%'", i+1)
			args = append(args, q)
		}
		conditions = append(conditions, strings.Join(queryConds, " AND "))
	}

	var codesParamIdx int
	if len(codes) > 0 {
		codesParamIdx = len(args) + 1
		conditions = append(conditions, fmt.Sprintf("code::text = ANY($%d::text[])", codesParamIdx))
		args = append(args, codes)
	}

	where := strings.Join(conditions, " AND ")

	var orderConds []string
	if len(queries) == 1 {
		orderConds = append(orderConds, fmt.Sprintf("similarity(search_text, $%d) DESC", firstQueryIdx))
	}
	if len(codes) > 0 {
		orderConds = append(orderConds, fmt.Sprintf("array_position($%d::text[], code::text)", codesParamIdx))
	}

	orderBy := ""
	if len(orderConds) > 0 {
		orderBy = "ORDER BY " + strings.Join(orderConds, ", ")
	}

	queryStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s %s", columns, Tables.Prices, where, orderBy)

	rows, err := r.db.Query(ctx, queryStr, args...)
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
	_, err := r.getExec(tx).CopyFrom(ctx, pgx.Identifier{Tables.Prices}, columns, pgx.CopyFromRows(rows))
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
			search_text = EXCLUDED.search_text`,
		Tables.Prices,
	)

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

	sql := fmt.Sprintf(`DELETE FROM %s WHERE code::text = ANY($1::text[])`, Tables.Prices)
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
	parts := []string{p.CurrentName, p.NewName, p.Template}
	text := strings.Join(parts, " ")
	return normalizeSearch(text)
}

func normalizeSearch(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "х", "x")
	return s
}
