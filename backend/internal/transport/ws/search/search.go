package search

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/ws/router"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
)

type WSHandler struct {
	service services.Search
}

func NewHandler(service services.Search) *WSHandler {
	return &WSHandler{
		service: service,
	}
}

func Register(router *router.WSRouter, service services.Search) {
	handler := NewHandler(service)

	router.Register("search", handler.Search)
	router.Register("search_stream", handler.SearchStream)
}

func (h *WSHandler) Search(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error {
	req := &models.SearchRequest{}
	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("failed to parse json. err: %w", err)
	}

	results, err := h.service.Search(ctx, req)
	if err != nil {
		return err
	}

	return client.SendJSON("SEARCH_RESULT", results)
}

func (h *WSHandler) SearchStream(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error {
	req := &models.SearchRequest{}
	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("failed to parse json. err: %w", err)
	}

	results, err := h.service.Search(ctx, req)
	if err != nil {
		return err
	}

	const batchSize = 50
	total := len(results)

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		// Формируем промежуточный payload
		payload := models.SearchResultPart{
			Items:  results[i:end],
			IsLast: end == total,
			Total:  total,
		}

		// Отправляем маленькую часть.
		// Внутри SendJSON должен стоять короткий SetWriteDeadline (напр. 5-10 сек)
		if err := client.SendJSON("SEARCH_RESULT_PART", payload); err != nil {
			return fmt.Errorf("failed to send batch %d: %w", i, err)
		}

		// Даем планировщику Go и сетевому стеку передохнуть
		// Это предотвращает блокировку event loop
		runtime.Gosched()
	}

	return nil
}
