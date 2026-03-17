package services

import (
	"context"
	"log"

	"github.com/Alexander272/Identic/backend/internal/repository"
)

type MessageBroadcaster interface {
	Broadcast(data []byte)
}

type OrderStreamService struct {
	repo repository.OrdersEvents
	hub  MessageBroadcaster
}

func NewOrderStreamService(repo repository.OrdersEvents, hub MessageBroadcaster) *OrderStreamService {
	return &OrderStreamService{
		repo: repo,
		hub:  hub,
	}
}

type OrdersStream interface {
	StartStreaming(ctx context.Context)
	Run(ctx context.Context) error
}

func (u *OrderStreamService) StartStreaming(ctx context.Context) {
	events := make(chan []byte)

	// Запускаем слушателя репозитория
	go u.repo.ListenOrders(ctx, events)

	// Читаем из канала и отправляем в хаб
	go func() {
		for msg := range events {
			u.hub.Broadcast(msg)
		}
	}()
}

func (u *OrderStreamService) Run(ctx context.Context) error {
	events := make(chan []byte)

	// Запускаем инфраструктурный слушатель.
	// Передаем контекст, чтобы репозиторий тоже знал об остановке.
	go u.repo.ListenOrders(ctx, events)

	log.Println("Order stream runner started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Order stream runner stopped")
			return ctx.Err() // Возвращаем причину остановки
		case data := <-events:
			// Ваша бизнес-логика трансляции в сокеты
			u.hub.Broadcast(data)
		}
	}
}
