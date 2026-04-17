package search_logs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/transport/ws/router"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
)

type WSHandler struct {
	service SearchLogGetter
}

type SearchLogGetter interface {
	Get(ctx context.Context, dto *models.GetSearchLogsDTO) ([]*models.SearchLog, error)
}

func NewHandler(service SearchLogGetter) *WSHandler {
	return &WSHandler{service: service}
}

func Register(router *router.WSRouter, service SearchLogGetter) {
	handler := NewHandler(service)
	router.Register("get_search_logs", handler.Get)
}

func (h *WSHandler) Get(ctx context.Context, client *ws_hub.Client, hub *ws_hub.Hub, data []byte) error {
	dto := &models.GetSearchLogsDTO{}
	if err := json.Unmarshal(data, dto); err != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}

	if dto.Limit == 0 {
		dto.Limit = 100
	}

	logs, err := h.service.Get(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to get search logs: %w", err)
	}

	return client.SendJSON("SEARCH_LOGS_RESULT", logs)
}
