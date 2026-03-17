package postgres

import (
	"context"
	"time"

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
			time.Sleep(5 * time.Second)
			continue
		}

		pgConn := conn.Conn()
		pgConn.Exec(ctx, "LISTEN order_updates")

		for {
			notification, err := pgConn.WaitForNotification(ctx)
			if err != nil {
				break // Ошибка соединения, выходим для реконнекта
			}
			// Передаем сырые байты в канал
			eventChan <- []byte(notification.Payload)
		}
		conn.Release()
	}
}
