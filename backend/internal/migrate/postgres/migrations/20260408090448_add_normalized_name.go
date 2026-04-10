package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddNormalizedName, downAddNormalizedName)
}

func upAddNormalizedName(ctx context.Context, tx *sql.Tx) error {
	// 1. Добавляем колонку
	if _, err := tx.ExecContext(ctx, `ALTER TABLE public.positions ADD COLUMN IF NOT EXISTS normalized_notes TEXT DEFAULT '';`); err != nil {
		return err
	}

	// 2. Читаем ВСЕ данные в память (если их до 100-200к, это нормально)
	rows, err := tx.QueryContext(ctx, "SELECT id, notes FROM public.positions WHERE notes != ''")
	if err != nil {
		return err
	}

	type data struct {
		id  string
		val string
	}
	var batch []data
	for rows.Next() {
		var d data
		if err := rows.Scan(&d.id, &d.val); err != nil {
			return err
		}
		batch = append(batch, data{id: d.id, val: services.NormalizeString(d.val)})
	}
	rows.Close() // Обязательно закрываем перед следующими запросами

	// 3. Обновляем через временную таблицу или UNNEST (очень быстро)
	// Это один тяжелый запрос вместо тысячи мелких
	if len(batch) > 0 {
		ids := make([]string, len(batch))
		vals := make([]string, len(batch))
		for i, b := range batch {
			ids[i] = b.id
			vals[i] = b.val
		}

		query := `
            UPDATE public.positions AS p
            SET normalized_notes = v.new_val
            FROM (SELECT unnest($1::text[]) as id, unnest($2::text[]) as new_val) as v
            WHERE p.id::text = v.id;`

		if _, err := tx.ExecContext(ctx, query, ids, vals); err != nil {
			return fmt.Errorf("failed to batch update: %w", err)
		}
	}

	return nil
}

func downAddNormalizedName(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "ALTER TABLE public.positions DROP COLUMN IF EXISTS normalized_notes;")
	return err
}
