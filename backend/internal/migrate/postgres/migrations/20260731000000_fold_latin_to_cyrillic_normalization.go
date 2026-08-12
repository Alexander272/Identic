package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFoldLatinToCyrillicNormalization, downFoldLatinToCyrillicNormalization)
}

func upFoldLatinToCyrillicNormalization(ctx context.Context, tx *sql.Tx) error {
	if err := refoldPositions(ctx, tx); err != nil {
		return err
	}

	return recalcRecentContentHashes(ctx, tx)
}

func refoldPositions(ctx context.Context, tx *sql.Tx) error {
	rows, err := tx.QueryContext(ctx, "SELECT id, name, search, COALESCE(notes, ''), normalized_notes FROM public.positions")
	if err != nil {
		return err
	}
	defer rows.Close()

	type fixedData struct {
		id              string
		search          string
		normalizedNotes string
	}
	var batch []fixedData
	for rows.Next() {
		var id, name, search, normalizedNotes, notes string
		if err := rows.Scan(&id, &name, &search, &notes, &normalizedNotes); err != nil {
			return err
		}

		wantSearch := services.NormalizeString(name)
		wantNormalizedNotes := services.NormalizeString(notes)

		if search != wantSearch || normalizedNotes != wantNormalizedNotes {
			batch = append(batch, fixedData{
				id:              id,
				search:          wantSearch,
				normalizedNotes: wantNormalizedNotes,
			})
		}
	}
	rows.Close()

	if len(batch) == 0 {
		return nil
	}

	ids := make([]string, len(batch))
	searches := make([]string, len(batch))
	normalizedNotes := make([]string, len(batch))
	for i, b := range batch {
		ids[i] = b.id
		searches[i] = b.search
		normalizedNotes[i] = b.normalizedNotes
	}

	query := `
		UPDATE public.positions AS p
		SET search = v.search, normalized_notes = v.normalized_notes
		FROM (
			SELECT
				unnest($1::text[]) as id,
				unnest($2::text[]) as search,
				unnest($3::text[]) as normalized_notes
		) as v
		WHERE p.id::text = v.id;`

	if _, err := tx.ExecContext(ctx, query, ids, searches, normalizedNotes); err != nil {
		return fmt.Errorf("failed to batch update positions: %w", err)
	}

	return nil
}

func recalcRecentContentHashes(ctx context.Context, tx *sql.Tx) error {
	var orderIDs []string
	rows, err := tx.QueryContext(ctx, `
		SELECT id FROM orders
		WHERE content_hash IS NOT NULL
		AND created_at > NOW() - INTERVAL '30 days'`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		orderIDs = append(orderIDs, id)
	}
	rows.Close()

	for _, orderID := range orderIDs {
		posRows, err := tx.QueryContext(ctx, `SELECT name, quantity FROM positions WHERE order_id = $1`, orderID)
		if err != nil {
			return err
		}

		var pts []*models.PositionDTO
		for posRows.Next() {
			p := &models.PositionDTO{}
			if err := posRows.Scan(&p.Name, &p.Quantity); err != nil {
				posRows.Close()
				return err
			}
			pts = append(pts, p)
		}
		posRows.Close()

		if len(pts) == 0 {
			continue
		}

		hash := services.CalculateHash(pts)

		if _, err := tx.ExecContext(ctx, `UPDATE orders SET content_hash = $1 WHERE id = $2`, hash, orderID); err != nil {
			return err
		}
	}

	return nil
}

func downFoldLatinToCyrillicNormalization(ctx context.Context, tx *sql.Tx) error {
	return nil
}
