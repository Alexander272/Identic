package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchRepo struct {
	db *pgxpool.Pool
}

func NewSearchRepo(db *pgxpool.Pool) *SearchRepo {
	return &SearchRepo{
		db: db,
	}
}

type Search interface {
	FetchExact(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error)
	FetchFuzzy(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error)
	FetchExactByQuantity(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error)
	FetchFuzzyByQuantity(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error)
}

func (r *SearchRepo) FetchExact(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	ids := make([]int, len(req.Items))
	names := make([]string, len(req.Items))
	qtys := make([]float64, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.Id
		names[i] = item.Name
		qtys[i] = item.Quantity
	}

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                id as req_item_id,
                name as req_name,
                qty as req_qty
            FROM UNNEST($1::text[], $2::numeric[], $3::int[]) AS t(name, qty, id)
        )
        SELECT 
            o.id::text, o.year, o.customer, o.consumer, o.date, o.is_bargaining, o.is_budget,
            r.req_item_id, p.id::text as matched_item_id, 
			CASE WHEN p.search=r.req_name THEN p.search ELSE p.normalized_notes END as p_search,
            r.req_qty, p.quantity as db_qty,
            1.0 as similarity -- Для точного поиска всегда 1.0
        FROM req r
        JOIN %s p ON p.search LIKE '%%' || r.req_name || '%%' OR p.normalized_notes LIKE '%%' || r.req_name || '%%'
        JOIN %s o ON o.id = p.order_id`,
		Tables.Positions, Tables.Orders,
	)

	rows, err := r.db.Query(ctx, query, names, qtys, ids)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.RawMatch
	for rows.Next() {
		var m models.RawMatch
		err := rows.Scan(
			&m.OrderId, &m.YearInt, &m.Customer, &m.Consumer, &m.Date, &m.IsBargaining, &m.IsBudget,
			&m.ReqId, &m.PosId, &m.PSearch, &m.ReqQty, &m.DbQty, &m.Similarity,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &m)
	}

	return results, nil
}

func (r *SearchRepo) FetchFuzzy(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	ids := make([]int, len(req.Items))
	names := make([]string, len(req.Items))
	qtys := make([]float64, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.Id
		names[i] = item.Name
		qtys[i] = item.Quantity
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, _ = tx.Exec(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.3;`)

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                id as req_id,
                name as req_name,
                qty as req_qty,
                regexp_split_to_array(name, '[^a-zA-Zа-яА-Я0-9.-]+') as req_tokens
            FROM UNNEST($1::text[], $2::numeric[], $3::int[]) AS t(name, qty, id)
        )
        SELECT 
            o.id::text, o.year, o.customer, o.consumer, o.date, o.is_bargaining, o.is_budget,
            r.req_id::text,
            p.id::text as pos_id,
            p.search,
            r.req_tokens,
            r.req_qty,
            p.quantity as db_qty,
            word_similarity(r.req_name, p.search) as sml
        FROM req r
        JOIN %s p ON p.search %% r.req_name 
        JOIN %s o ON o.id = p.order_id;`,
		Tables.Positions, Tables.Orders,
	)

	rows, err := tx.Query(ctx, query, names, qtys, ids)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.RawMatch
	for rows.Next() {
		var m models.RawMatch
		// Порядок Scan должен строго соответствовать порядку колонок в SELECT
		err := rows.Scan(
			&m.OrderId,
			&m.YearInt,
			&m.Customer,
			&m.Consumer,
			&m.Date,
			&m.IsBargaining,
			&m.IsBudget,
			&m.ReqId,
			&m.PosId,
			&m.PSearch,
			&m.ReqTokens,
			&m.ReqQty,
			&m.DbQty,
			&m.Similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results, &m)
	}

	return results, nil
}

func (r *SearchRepo) FetchExactByQuantity(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	ids := make([]int, len(req.Items))
	qtys := make([]float64, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.Id
		qtys[i] = item.Quantity
	}

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                id as req_item_id,
                qty as req_qty
            FROM UNNEST($1::numeric[], $2::int[]) AS t(qty, id)
        )
        SELECT 
            o.id::text, o.year, o.customer, o.consumer, o.date, o.is_bargaining, o.is_budget,
            r.req_item_id, p.id::text as matched_item_id, 
			p.search as p_search,
            r.req_qty, p.quantity as db_qty,
            1.0 as similarity
        FROM req r
        JOIN %s p ON p.quantity = r.req_qty
        JOIN %s o ON o.id = p.order_id`,
		Tables.Positions, Tables.Orders,
	)

	rows, err := r.db.Query(ctx, query, qtys, ids)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.RawMatch
	for rows.Next() {
		var m models.RawMatch
		err := rows.Scan(
			&m.OrderId, &m.YearInt, &m.Customer, &m.Consumer, &m.Date, &m.IsBargaining, &m.IsBudget,
			&m.ReqId, &m.PosId, &m.PSearch, &m.ReqQty, &m.DbQty, &m.Similarity,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &m)
	}

	return results, nil
}

func (r *SearchRepo) FetchFuzzyByQuantity(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	ids := make([]int, len(req.Items))
	qtys := make([]float64, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.Id
		qtys[i] = item.Quantity
	}

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                id as req_item_id,
                qty as req_qty
            FROM UNNEST($1::numeric[], $2::int[]) AS t(qty, id)
        )
        SELECT 
            o.id::text, o.year, o.customer, o.consumer, o.date, o.is_bargaining, o.is_budget,
            r.req_item_id, p.id::text as matched_item_id, 
			p.search as p_search,
            r.req_qty, p.quantity as db_qty,
            LEAST(p.quantity, r.req_qty) / GREATEST(p.quantity, r.req_qty) as similarity
        FROM req r
        JOIN %s p ON p.quantity BETWEEN r.req_qty * 0.7 AND r.req_qty * 1.3
        JOIN %s o ON o.id = p.order_id`,
		Tables.Positions, Tables.Orders,
	)

	rows, err := r.db.Query(ctx, query, qtys, ids)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.RawMatch
	for rows.Next() {
		var m models.RawMatch
		err := rows.Scan(
			&m.OrderId, &m.YearInt, &m.Customer, &m.Consumer, &m.Date, &m.IsBargaining, &m.IsBudget,
			&m.ReqId, &m.PosId, &m.PSearch, &m.ReqQty, &m.DbQty, &m.Similarity,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &m)
	}

	return results, nil
}
