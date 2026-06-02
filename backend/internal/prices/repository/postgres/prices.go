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

type PricesRepo struct {
	db *pgxpool.Pool
}

func NewPricesRepo(db *pgxpool.Pool) *PricesRepo {
	return &PricesRepo{db: db}
}

type Prices interface {
	GetAll(ctx context.Context, page, perPage int) ([]*models.Price, int, error)
	Search(ctx context.Context, queries, codes, fields []string, page, perPage int) ([]*models.Price, int, error)
	SearchAll(ctx context.Context, queries []string, codes []string) ([]*models.Price, error)
	CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.Price) error
	UpsertSeveral(ctx context.Context, tx postgres.Tx, dto []*models.Price) error
	DeleteSeveral(ctx context.Context, tx postgres.Tx, codes []string) error
}

func (r *PricesRepo) GetAll(ctx context.Context, page, perPage int) ([]*models.Price, int, error) {
	columns := "prices.id, prices.code, prices.current_name, prices.new_name, prices.price, prices.template, prices.note, prices.need_sibur_approval"

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", Tables.Prices)
	var total int
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, postgres.MapError(fmt.Errorf("failed to count: %w", err))
	}

	offset := (page - 1) * perPage
	dataQuery := fmt.Sprintf("SELECT %s FROM %s ORDER BY prices.code LIMIT $1 OFFSET $2", columns, Tables.Prices)

	rows, err := r.db.Query(ctx, dataQuery, perPage, offset)
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

func (r *PricesRepo) getExec(tx postgres.Tx) postgres.QueryExecutor {
	if tx != nil {
		return tx.TX()
	}
	return r.db
}

type fieldMapping struct {
	column string
	exact  bool
}

var allowedFields = map[string]fieldMapping{
	"current_name":        {"current_name_norm", false},
	"new_name":            {"new_name_norm", false},
	"template":            {"template_norm", false},
	"price":               {"price", true},
	"note":                {"note", false},
	"need_sibur_approval": {"need_sibur_approval", false},
}

func (r *PricesRepo) Search(ctx context.Context, queries, codes, fields []string, page, perPage int) ([]*models.Price, int, error) {
	columns := "prices.id, prices.code, prices.current_name, prices.new_name, prices.price, prices.template, prices.note, prices.need_sibur_approval"

	var conditions []string
	args := make([]any, 0)

	var firstQueryIdx int
	if len(queries) > 0 {
		firstQueryIdx = len(args) + 1

		if len(fields) > 0 {
			fieldMappers := make([]fieldMapping, 0, len(fields))
			for _, f := range fields {
				if fm, ok := allowedFields[f]; ok {
					fieldMappers = append(fieldMappers, fm)
				}
			}
			if len(fieldMappers) == 0 {
				fieldMappers = []fieldMapping{{column: "search_text", exact: false}}
			}

			queryConds := make([]string, len(queries))
			for i, q := range queries {
				fieldConds := make([]string, len(fieldMappers))
				for j, fm := range fieldMappers {
					if fm.exact {
						if _, err := strconv.ParseFloat(q, 64); err == nil {
							fieldConds[j] = fmt.Sprintf("%s = $%d::numeric", fm.column, i+1)
						} else {
							fieldConds[j] = "false"
						}
					} else {
						fieldConds[j] = fmt.Sprintf("%s ILIKE '%%' || $%d || '%%'", fm.column, i+1)
					}
				}
				queryConds[i] = "(" + strings.Join(fieldConds, " OR ") + ")"
				args = append(args, q)
			}
			conditions = append(conditions, strings.Join(queryConds, " AND "))
		} else {
			queryConds := make([]string, len(queries))
			for i, q := range queries {
				queryConds[i] = fmt.Sprintf("search_text ILIKE '%%' || $%d || '%%'", i+1)
				args = append(args, q)
			}
			conditions = append(conditions, strings.Join(queryConds, " AND "))
		}
	}

	// Старый вариант: code::text = ANY($N) + array_position — НЕ поддерживал дубликаты кодов и их порядок в запросе
	// var codesParamIdx int
	// if len(codes) > 0 {
	// 	codesParamIdx = len(args) + 1
	// 	conditions = append(conditions, fmt.Sprintf("code::text = ANY($%d::text[])", codesParamIdx))
	// 	args = append(args, codes)
	// }

	// JOIN UNNEST(...) WITH ORDINALITY — для каждого элемента массива создаётся строка,
	// дубликаты сохраняются, c.ord даёт порядок из запроса.
	var fromClause string
	if len(codes) > 0 {
		codesParamIdx := len(args) + 1
		fromClause = fmt.Sprintf("%s JOIN UNNEST($%d::text[]) WITH ORDINALITY AS c(code, ord) ON %s.code::text = c.code",
			Tables.Prices, codesParamIdx, Tables.Prices)
		args = append(args, codes)
	} else {
		fromClause = Tables.Prices
	}

	where := strings.Join(conditions, " AND ")

	whereClause := ""
	if where != "" {
		whereClause = "WHERE " + where
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", fromClause, whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, postgres.MapError(fmt.Errorf("failed to count: %w", err))
	}

	// var orderConds []string
	// if len(fields) == 0 && len(queries) == 1 {
	// 	orderConds = append(orderConds, fmt.Sprintf("similarity(search_text, $%d) DESC", firstQueryIdx))
	// }
	// if len(codes) > 0 {
	// 	orderConds = append(orderConds, fmt.Sprintf("array_position($%d::text[], code::text)", codesParamIdx))
	// }

	var orderConds []string
	if len(codes) > 0 {
		orderConds = append(orderConds, "c.ord")
	}
	if len(fields) == 0 && len(queries) == 1 {
		orderConds = append(orderConds, fmt.Sprintf("similarity(search_text, $%d) DESC", firstQueryIdx))
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

	dataQuery := fmt.Sprintf("SELECT %s FROM %s %s %s LIMIT $%d OFFSET $%d",
		columns, fromClause, whereClause, orderBy, perPageIdx, offsetIdx)

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
	columns := "prices.id, prices.code, prices.current_name, prices.new_name, prices.price, prices.template, prices.note, prices.need_sibur_approval"

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

	// Старый вариант: code::text = ANY($N) + array_position — НЕ поддерживал дубликаты кодов и их порядок в запросе
	// var codesParamIdx int
	// if len(codes) > 0 {
	// 	codesParamIdx = len(args) + 1
	// 	conditions = append(conditions, fmt.Sprintf("code::text = ANY($%d::text[])", codesParamIdx))
	// 	args = append(args, codes)
	// }

	// JOIN UNNEST(...) WITH ORDINALITY — для каждого элемента массива создаётся строка,
	// дубликаты сохраняются, c.ord даёт порядок из запроса.
	var fromClause string
	if len(codes) > 0 {
		codesParamIdx := len(args) + 1
		fromClause = fmt.Sprintf("%s JOIN UNNEST($%d::text[]) WITH ORDINALITY AS c(code, ord) ON %s.code::text = c.code",
			Tables.Prices, codesParamIdx, Tables.Prices)
		args = append(args, codes)
	} else {
		fromClause = Tables.Prices
	}

	where := strings.Join(conditions, " AND ")

	// var orderConds []string
	// if len(queries) == 1 {
	// 	orderConds = append(orderConds, fmt.Sprintf("similarity(search_text, $%d) DESC", firstQueryIdx))
	// }
	// if len(codes) > 0 {
	// 	orderConds = append(orderConds, fmt.Sprintf("array_position($%d::text[], code::text)", codesParamIdx))
	// }

	var orderConds []string
	if len(codes) > 0 {
		orderConds = append(orderConds, "c.ord")
	}
	if len(queries) == 1 {
		orderConds = append(orderConds, fmt.Sprintf("similarity(search_text, $%d) DESC", firstQueryIdx))
	}

	orderBy := ""
	if len(orderConds) > 0 {
		orderBy = "ORDER BY " + strings.Join(orderConds, ", ")
	}

	queryStr := "SELECT " + columns + " FROM " + fromClause
	if where != "" {
		queryStr += " WHERE " + where
	}
	if orderBy != "" {
		queryStr += " " + orderBy
	}

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
		rows[i] = []interface{}{
			p.Code, p.CurrentName, p.NewName, p.Price, p.Template, p.Note, p.NeedSiburApproval,
			p.SearchText,
			p.CurrentNameNorm, p.NewNameNorm, p.TemplateNorm,
		}
	}

	columns := []string{
		"code", "current_name", "new_name", "price", "template", "note", "need_sibur_approval", "search_text",
		"current_name_norm", "new_name_norm", "template_norm",
	}
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
		INSERT INTO %s (code, current_name, new_name, price, template, note, need_sibur_approval, search_text,
			current_name_norm, new_name_norm, template_norm)
		SELECT * FROM UNNEST($1::text[], $2::text[], $3::text[], $4::float8[], $5::text[], $6::text[], $7::text[], $8::text[],
			$9::text[], $10::text[], $11::text[])
		ON CONFLICT (code) DO UPDATE SET
			current_name = EXCLUDED.current_name,
			new_name = EXCLUDED.new_name,
			price = EXCLUDED.price,
			template = EXCLUDED.template,
			note = EXCLUDED.note,
			need_sibur_approval = EXCLUDED.need_sibur_approval,
			search_text = EXCLUDED.search_text,
			current_name_norm = EXCLUDED.current_name_norm,
			new_name_norm = EXCLUDED.new_name_norm,
			template_norm = EXCLUDED.template_norm`,
		Tables.Prices,
	)

	codes := make([]string, len(dto))
	currentNames := make([]string, len(dto))
	newNames := make([]string, len(dto))
	prices := make([]float64, len(dto))
	templates := make([]string, len(dto))
	notes := make([]string, len(dto))
	needSiburApprovals := make([]string, len(dto))
	searchTexts := make([]string, len(dto))
	currentNameNorms := make([]string, len(dto))
	newNameNorms := make([]string, len(dto))
	templateNorms := make([]string, len(dto))

	for i, p := range dto {
		codes[i] = p.Code
		currentNames[i] = p.CurrentName
		newNames[i] = p.NewName
		prices[i] = p.Price
		templates[i] = p.Template
		notes[i] = p.Note
		needSiburApprovals[i] = p.NeedSiburApproval
		searchTexts[i] = p.SearchText
		currentNameNorms[i] = p.CurrentNameNorm
		newNameNorms[i] = p.NewNameNorm
		templateNorms[i] = p.TemplateNorm
	}

	_, err := r.getExec(tx).Exec(ctx, sql,
		codes, currentNames, newNames, prices, templates, notes, needSiburApprovals, searchTexts,
		currentNameNorms, newNameNorms, templateNorms,
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
		if err := rows.Scan(&p.ID, &p.Code, &p.CurrentName, &p.NewName, &p.Price, &p.Template, &p.Note, &p.NeedSiburApproval); err != nil {
			return nil, postgres.MapError(fmt.Errorf("scan row error: %w", err))
		}
		data = append(data, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, postgres.MapError(fmt.Errorf("rows iteration error: %w", err))
	}
	return data, nil
}
