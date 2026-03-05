package repository

import (
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction interface {
	postgres.Transaction
}

type Orders interface {
	postgres.Orders
}
type Positions interface {
	postgres.Positions
}

type Repository struct {
	Transaction
	Orders
	Positions
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	transaction := postgres.NewTransactionRepo(pool)

	return &Repository{
		Transaction: transaction,
		Orders:      postgres.NewOrderRepo(pool, transaction),
		Positions:   postgres.NewPositionRepo(pool, transaction),
	}
}
