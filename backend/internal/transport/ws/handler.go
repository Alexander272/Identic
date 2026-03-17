package ws

import (
	"log"
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/gorilla/websocket"
)

func NewUpgrader(allowedOrigins []string) websocket.Upgrader {
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

func HandleWS(hub *Hub, conf config.HttpConfig, w http.ResponseWriter, r *http.Request) {
	upgrader := NewUpgrader(conf.AllowedOrigins)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	// Запускаем два процесса для каждого клиента:

	// 1. Запись (из Go в Браузер)
	go func() {
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()
		for message := range client.send {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				break
			}
		}
	}()

	// 2. Чтение (из Браузера в Go) — нужно, чтобы вовремя заметить разрыв связи
	go func() {
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()
		for {
			// Мы не ожидаем данных от клиента в этой задаче,
			// но чтение необходимо для обработки Ping/Pong и Close
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}()
}
