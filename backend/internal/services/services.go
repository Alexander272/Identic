package services

import (
	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
)

type MessageBroadcaster interface {
	BroadcastMessage(topic string, data []byte)
}

type Services struct {
	Import
	Orders
	OrdersStream
	Positions
	Search
	SearchStream
}

type Deps struct {
	Repo  *repository.Repository
	Links config.LinksConfig
	Hub   *ws_hub.Hub
}

func NewServices(deps *Deps) *Services {
	transaction := NewTransactionManager(deps.Repo.Transaction)

	positions := NewPositionsService(deps.Repo.Positions, transaction)
	orders := NewOrdersService(deps.Repo.Orders, transaction, positions)
	import_file := NewImportService(transaction, orders, positions)

	ordersStream := NewOrderStreamService(deps.Repo.OrdersEvents, deps.Hub)

	search := NewSearchService(deps.Repo.Search, deps.Links.Orders)
	searchStream := NewSearchStreamService(search, deps.Hub)

	return &Services{
		Import:       import_file,
		Orders:       orders,
		OrdersStream: ordersStream,
		Positions:    positions,
		Search:       search,
		SearchStream: searchStream,
	}
}
