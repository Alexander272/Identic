package router

import (
	"context"
	"strings"
	"sync"

	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/goccy/go-json"
)

type WSHandler func(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error

type WSRouter struct {
	handlers map[string]WSHandler
	mu       sync.RWMutex
}

func NewWSRouter() *WSRouter {
	return &WSRouter{
		handlers: make(map[string]WSHandler),
	}
}

// Register регистрирует обработчик для конкретного действия (action)
func (r *WSRouter) Register(action string, h WSHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[action] = h
}

// Handle маршрутизирует входящее сообщение к нужному обработчику
func (r *WSRouter) Handle(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, rawMsg []byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Panic in WS Handler",
				logger.AnyAttr("recovery", err),
				logger.StringAttr("client_addr", client.Conn.RemoteAddr().String()),
			)
			client.SendJSON("error", map[string]string{"message": "Internal server error"})
		}
	}()

	// 1. Парсим общую обертку команды
	var cmd struct {
		Action  string          `json:"action"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := json.Unmarshal(rawMsg, &cmd); err != nil {
		logger.Info("WS JSON parse error", logger.ErrAttr(err))
		client.SendJSON("error", map[string]string{"message": "Invalid JSON format"})
		return
	}

	if cmd.Action == "" {
		client.SendJSON("error", map[string]string{"message": "Action is required"})
		return
	}

	// 2. Ищем обработчик
	action := strings.ToLower(cmd.Action)

	r.mu.RLock()
	handler, exists := r.handlers[action]
	r.mu.RUnlock()

	if !exists {
		logger.Info("Unknown action:", logger.StringAttr("action", cmd.Action))
		client.SendJSON("error", map[string]string{"message": "Unknown action: " + cmd.Action})
		return
	}

	// 3. Вызываем обработчик
	// Обработка ошибок внутри самого хендлера или здесь
	if err := handler(ctx, client, hub, cmd.Payload); err != nil {
		// Глобальная обработка ошибок, если хендлер вернул ошибку, но не отправил ответ сам
		logger.Error("Handler error", logger.StringAttr("action", cmd.Action), logger.ErrAttr(err))
		// Можно отправить универсальную ошибку, если хендлер этого не сделал
		client.SendJSON("error", map[string]string{"message": "Произошла ошибка: " + err.Error()})
	}
}
