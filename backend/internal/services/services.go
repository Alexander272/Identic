package services

import "github.com/Alexander272/Identic/backend/internal/repository"

type Services struct {
	Import
	Orders
	Positions
}

type Deps struct {
	Repo *repository.Repository
}

func NewServices(deps *Deps) *Services {
	transaction := NewTransactionManager(deps.Repo.Transaction)

	positions := NewPositionsService(deps.Repo.Positions, transaction)
	orders := NewOrdersService(deps.Repo.Orders, transaction, positions)
	import_file := NewImportService(transaction, orders, positions)

	return &Services{
		Import:    import_file,
		Orders:    orders,
		Positions: positions,
	}
}
