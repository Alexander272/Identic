package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upBackfillSearchAndNormalizedNotes, downBackfillSearchAndNormalizedNotes)
}

func upBackfillSearchAndNormalizedNotes(ctx context.Context, tx *sql.Tx) error {
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
		return fmt.Errorf("failed to batch update: %w", err)
	}

	return nil
}

func downBackfillSearchAndNormalizedNotes(ctx context.Context, tx *sql.Tx) error {
	return nil
}
