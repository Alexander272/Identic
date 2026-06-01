package repository

import (
	"context"

	pg "github.com/Alexander272/Identic/backend/internal/prices/repository/postgres"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(tx postgres.Tx) error) error
}

type Prices interface {
	pg.Prices
}

type PriceSearchLogs interface {
	pg.PriceSearchLogs
}

type Repository struct {
	Prices
	PriceSearchLogs
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Prices:          pg.NewPricesRepo(db),
		PriceSearchLogs: pg.NewPriceSearchLogRepo(db),
	}
}
