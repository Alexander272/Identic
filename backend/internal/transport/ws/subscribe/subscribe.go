package subscribe

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/transport/ws/router"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/goccy/go-json"
)

type WSHandler struct{}

func NewHandler() *WSHandler {
	return &WSHandler{}
}

func Register(router *router.WSRouter) {
	handler := NewHandler()

	router.Register("subscribe", handler.Subscribe)
	router.Register("unsubscribe", handler.Unsubscribe)
}

func (h *WSHandler) Subscribe(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error {
	topic := &models.Subscribe{}
	// Предполагаем, что data — это просто строка с названием топика в JSON
	if err := json.Unmarshal(data, &topic); err != nil {
		return fmt.Errorf("failed to unmarshal topic: %w", err)
	}

	// Создаем объект подписки
	sub := &ws_hub.Subscription{
		Client: client,
		Topic:  topic.Topic,
	}

	select {
	case hub.Register <- sub:
		// Успешно отправили запрос в очередь хаба
		return client.SendJSON("SUBSCRIBED", map[string]string{"topic": topic.Topic})
	default:
		// Канал хаба переполнен (busy)
		logger.Error("Hub register channel full", logger.StringAttr("topic", topic.Topic))
		return client.SendJSON("ERROR", "server_busy")
	}
}

func (h *WSHandler) Unsubscribe(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error {
	topic := &models.Subscribe{}
	if err := json.Unmarshal(data, &topic); err != nil {
		return fmt.Errorf("failed to unmarshal topic: %w", err)
	}

	select {
	case hub.Unregister <- &ws_hub.Subscription{Client: client, Topic: topic.Topic}:
		// Обычно подтверждение отписки не требуется, но можно отправить
		return client.SendJSON("UNSUBSCRIBED", map[string]string{"topic": topic.Topic})
	default:
		logger.Info("Hub unregister channel full", logger.StringAttr("topic", topic.Topic))
	}
	return nil
}

// func (h *WSHandler) Subscribe(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error {
// 	var topic string
// 	if err := json.Unmarshal(data, &topic); err != nil {
// 		return err
// 	}

// 	req := &ws_hub.SubscriptionRequest{Client: client, Topic: topic}
// 	select {
// 	case hub.Register <- req:
// 		client.SendJSON("subscribed", map[string]string{"topic": topic})
// 	default:
// 		client.SendJSON("error", map[string]string{"message": "Server busy"})
// 	}
// 	return nil
// }

// func (h *WSHandler) Unsubscribe(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error {
// 	var topic string
// 	if err := json.Unmarshal(data, &topic); err != nil {
// 		return err
// 	}

// 	req := &ws_hub.UnsubscriptionRequest{Client: client, Topic: topic}
// 	select {
// 	case hub.Unregister <- req:
// 	default:
// 		logger.Info("Hub unregister channel full")
// 	}
// 	return nil
// }
