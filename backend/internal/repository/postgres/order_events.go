package postgres

import (
	"context"
	"time"

	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderEventRepo struct {
	pool *pgxpool.Pool
}

func NewOrderEventRepo(pool *pgxpool.Pool) *OrderEventRepo {
	return &OrderEventRepo{
		pool: pool,
	}
}

type OrderEvent interface {
	ListenOrders(ctx context.Context, eventChan chan<- []byte)
}

func (r *OrderEventRepo) ListenOrders(ctx context.Context, eventChan chan<- []byte) {
	for {
		conn, err := r.pool.Acquire(ctx)
		if err != nil {
			logger.Error("failed to acquire connection", logger.ErrAttr(err))
			time.Sleep(5 * time.Second)
			continue
		}

		pgConn := conn.Conn()
		_, err = pgConn.Exec(ctx, "LISTEN order_updates")
		if err != nil {
			logger.Error("Не удалось выполнить LISTEN", logger.ErrAttr(err))
			conn.Release()
			time.Sleep(5 * time.Second)
			continue
		}

		for {
			notification, err := pgConn.WaitForNotification(ctx)
			if err != nil {
				logger.Error("Соединение разорвано или контекст отменен", logger.ErrAttr(err))
				break // Ошибка соединения, выходим для реконнекта
			}
			// Передаем сырые байты в канал
			eventChan <- []byte(notification.Payload)
		}
		conn.Release()
	}
}
