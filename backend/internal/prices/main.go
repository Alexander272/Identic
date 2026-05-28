package prices

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Alexander272/Identic/backend/internal/prices/repository"
	"github.com/Alexander272/Identic/backend/internal/prices/services"
	transport "github.com/Alexander272/Identic/backend/internal/prices/transport/http"
)

func NewPricesModule(db *pgxpool.Pool, tm repository.TransactionManager) *transport.Handler {
	repo := repository.NewRepository(db)
	svc := services.NewServices(&services.Deps{Repos: repo, TxManager: tm})
	handler := transport.NewHandler(svc)
	return handler
}
