package ws

import (
	"context"
	"net/http"
	"time"

	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/ws/router"
	"github.com/Alexander272/Identic/backend/internal/transport/ws/search"
	"github.com/Alexander272/Identic/backend/internal/transport/ws/search_logs"
	"github.com/Alexander272/Identic/backend/internal/transport/ws/subscribe"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	hub      *ws_hub.Hub
	conf     config.HttpConfig
	services *services.Services
	router   *router.WSRouter
}

func NewWsHandler(hub *ws_hub.Hub, conf config.HttpConfig, services *services.Services) *WsHandler {
	router := router.NewWSRouter()

	wsHandler := &WsHandler{
		hub:      hub,
		conf:     conf,
		services: services,
		router:   router,
	}

	router.Register("PING", wsHandler.HandlePing)
	subscribe.Register(router)
	search.Register(router, wsHandler.services.Search)
	search_logs.Register(router, wsHandler.services.SearchStream.GetSearchLogService())

	return wsHandler
}

func (h *WsHandler) NewUpgrader(allowedOrigins []string) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")

			// Если список пуст, разрешаем только локальные запросы (безопасный дефолт)
			if len(allowedOrigins) == 0 {
				return origin == ""
			}

			for _, o := range allowedOrigins {
				if origin == o {
					return true
				}
			}
			return false
		},
	}
}

func (h *WsHandler) HandleWS(c *gin.Context) {
	upgrader := h.NewUpgrader(h.conf.AllowedOrigins)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Upgrade error", logger.StringAttr("error", err.Error()))
		return
	}

	client := ws_hub.NewClient(conn, h.hub)

	go h.clientWriter(client)
	go h.clientReader(client)
}

// clientWriter - отвечает ТОЛЬКО за отправку данных в сокет.
// Не знает ни о каких командах, подписках или поиске.
func (h *WsHandler) clientWriter(client *ws_hub.Client) {
	pingPeriod := h.conf.PingPeriod // Как часто слать пинги
	writeWait := h.conf.WriteWait   // Таймаут на саму запись в сокет

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			// Устанавливаем дедлайн на запись, чтобы не "зависнуть" на мертвом сокете
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// Хаб закрыл канал (значит клиент удален)
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				client.Hub.Disconnect(client)
				return
			}

		case <-ticker.C:
			// Каждые 45 секунд шлем Ping
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				client.Hub.Disconnect(client)
				return
			}
		}
	}
}

// clientReader - отвечает ТОЛЬКО за чтение байтов из сокета и передачу их диспетчеру.
func (h *WsHandler) clientReader(client *ws_hub.Client) {
	pongWait := h.conf.PongWait // Время ожидания ответа от клиента

	defer client.Hub.Disconnect(client)

	client.Conn.SetReadLimit(h.conf.MaxMessageSize)

	// 1. Устанавливаем начальный дедлайн на чтение
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// 2. Настраиваем обработчик Pong-сообщений
	// Каждый раз, когда клиент отвечает на наш Ping, мы "отодвигаем" дедлайн дальше
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNormalClosure,   // 1000 - штатное закрытие
				websocket.CloseGoingAway,       // 1001 - уход со страницы
				websocket.CloseAbnormalClosure, // 1006 - обрыв связи
			) {
				logger.Error("Unexpected close error:", logger.StringAttr("error", err.Error()))
			}
			return
		}

		// Создаем контекст для КОНКРЕТНОГО сообщения
		// Используем Background, так как сокет - это долгоживущее соединение
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

		// Запускаем обработку в горутине, чтобы ReadMessage не блокировался
		// и мог продолжать принимать Pong-пакеты, пока идет "тяжелый" поиск
		go func(c context.Context, can context.CancelFunc, msg []byte) {
			defer can()
			h.router.Handle(c, client, h.hub, msg)
		}(ctx, cancel, rawMsg)
	}
}

// func (h *WsHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
// 	h.router.Register("PING", h.HandlePing)
// 	search.Register(h.router, h.services.Search)

// 	upgrader := h.NewUpgrader(h.conf.AllowedOrigins)

// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("Upgrade error: %v", err)
// 		return
// 	}

// 	client := &ws_hub.Client{Conn: conn, Send: make(chan []byte, 256)}
// 	h.hub.Register <- client

// 	// Запускаем два процесса для каждого клиента:

// 	// // 1. Запись (из Go в Браузер)
// 	// go client.WritePump()
// 	// // 2. Чтение (из Браузера в Go) — нужно, чтобы вовремя заметить разрыв связи
// 	// go client.ReadPump()

// 	// 1. Запись (из Go в Браузер)
// 	go func() {
// 		defer func() {
// 			h.hub.Unregister <- client
// 			conn.Close()
// 		}()
// 		for message := range client.Send {
// 			err := conn.WriteMessage(websocket.TextMessage, message)
// 			if err != nil {
// 				break
// 			}
// 		}
// 	}()

// 	// 2. Чтение (из Браузера в Go) — нужно, чтобы вовремя заметить разрыв связи
// 	go func() {
// 		defer func() {
// 			h.hub.Unregister <- client
// 			conn.Close()
// 		}()
// 		for {
// 			// Мы не ожидаем данных от клиента в этой задаче,
// 			// но чтение необходимо для обработки Ping/Pong и Close
// 			if _, _, err := conn.NextReader(); err != nil {
// 				break
// 			}
// 		}
// 	}()
// }

func (h *WsHandler) HandlePing(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error {
	return client.SendJSON("PONG", nil)
}
