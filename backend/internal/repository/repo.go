package repository

import (
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/internal/repository/redis"
	memoryDB "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction interface {
	postgres.Transaction
}

type Orders interface {
	postgres.Orders
}
type OrdersEvents interface {
	postgres.OrderEvent
}
type Positions interface {
	postgres.Positions
}

type Permissions interface {
	postgres.Permissions
}
type Roles interface {
	postgres.Roles
}
type RoleHierarchy interface {
	postgres.RoleHierarchy
}
type Users interface {
	postgres.Users
}

// type Search struct {
// 	pgRepo    postgres.Search // конкретный репозиторий Postgres
// 	redisRepo redis.Search    // конкретный репозиторий Redis
// }

// type Search interface {
// 	searchAggregator
// }

type Repository struct {
	Transaction
	Orders
	OrdersEvents
	Positions
	Search

	Permissions
	Roles
	RoleHierarchy
	Users
}

func NewRepository(pool *pgxpool.Pool, memDB *memoryDB.Client) *Repository {
	transaction := postgres.NewTransactionRepo(pool)

	return &Repository{
		Transaction:  transaction,
		OrdersEvents: postgres.NewOrderEventRepo(pool),
		Orders:       postgres.NewOrderRepo(pool, transaction),
		Positions:    postgres.NewPositionRepo(pool, transaction),
		Search: searchProvider{
			postgresRepo: postgres.NewSearchRepo(pool),
			redisRepo:    redis.NewSearchRepo(memDB),
		},

		Permissions:   postgres.NewPermissionRepo(pool, transaction),
		Roles:         postgres.NewRoleRepo(pool, transaction),
		RoleHierarchy: postgres.NewRoleHierarchyRepo(pool, transaction),
		Users:         postgres.NewUserRepo(pool, transaction),
	}
}
