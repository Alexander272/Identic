package migrations

import (
	"context"
	"database/sql"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFillContentHash, downFillContentHash)
}

func upFillContentHash(ctx context.Context, tx *sql.Tx) error {
	var orderIDs []string
	rows, err := tx.QueryContext(ctx, `
		SELECT id FROM orders 
		WHERE content_hash IS NULL 
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

	// 2. Обрабатываем каждый заказ отдельно
	for _, orderID := range orderIDs {
		// Достаем позиции
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

		// Считаем хеш через общий пакет
		hash := services.CalculateHash(pts)

		if _, err := tx.ExecContext(ctx, `UPDATE orders SET content_hash = $1 WHERE id = $2`, hash, orderID); err != nil {
			return err
		}
	}

	return nil
}

func downFillContentHash(ctx context.Context, tx *sql.Tx) error {
	// При откате просто очищаем колонку (индекс и сама колонка удалятся SQL-файлом)
	_, err := tx.ExecContext(ctx, "UPDATE orders SET content_hash = NULL")
	return err
}
