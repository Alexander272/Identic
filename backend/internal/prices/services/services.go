package services

import (
	"context"

	"github.com/Alexander272/Identic/backend/internal/prices/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
)

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(tx postgres.Tx) error) error
}

type Services struct {
	Prices
	Export
	Import
	PriceSearchLogs
}

type Deps struct {
	Repos     *repository.Repository
	TxManager repository.TransactionManager
}

func NewServices(deps *Deps) *Services {
	searchLog := NewPriceSearchLogService(deps.Repos.PriceSearchLogs)
	prices := NewPricesService(deps.Repos.Prices, deps.TxManager, searchLog)
	export := NewExportService(deps.Repos.Prices)
	importSvc := NewImportService(deps.Repos.Prices, deps.TxManager)

	return &Services{
		Prices:          prices,
		Export:          export,
		Import:          importSvc,
		PriceSearchLogs: searchLog,
	}
}
