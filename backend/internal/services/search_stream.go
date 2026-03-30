package services

import (
	"context"
	"time"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/goccy/go-json"
)

type SearchStreamService struct {
	service Search
	hub     MessageBroadcaster
}

func NewSearchStreamService(service Search, hub MessageBroadcaster) *SearchStreamService {
	return &SearchStreamService{
		service: service,
		hub:     hub,
	}
}

type SearchStream interface {
	Streaming(ctx context.Context, req *models.SearchRequest)
}

func (s *SearchStreamService) Streaming(ctx context.Context, req *models.SearchRequest) {
	start := time.Now()
	results, err := s.service.Search(ctx, req)
	if err != nil {
		s.sendError(req.SearchId, err)
		return
	}
	logger.Debug("search",
		logger.StringAttr("search_id", req.SearchId),
		logger.IntAttr("count", len(results)),
		logger.AnyAttr("time", time.Since(start)),
	)

	const batchSize = 10
	total := len(results)

	batch := make([]*models.OrderMatchResult, 0, batchSize)

	for _, item := range results {
		select {
		case <-ctx.Done():
			return
		default:
		}

		batch = append(batch, item)

		if len(batch) == batchSize {
			s.sendBatch(req.SearchId, batch, false, total)
			batch = make([]*models.OrderMatchResult, 0, batchSize)
		}
	}

	// отправляем остаток
	if len(batch) > 0 {
		s.sendBatch(req.SearchId, batch, false, total)
	}

	// финальное сообщение
	s.sendBatch(req.SearchId, []*models.OrderMatchResult{}, true, total)
}

func (s *SearchStreamService) sendBatch(searchId string, items []*models.OrderMatchResult, isLast bool, total int) {
	payload := models.SearchResultPart{
		Items:  items,
		IsLast: isLast,
		Total:  total,
	}

	msg := ws_hub.WSMessage{
		Action: "SEARCH_RESULT_PART",
		Data:   payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		s.sendError(searchId, err)
		return
	}

	s.sendPart(searchId, data)
}

func (s *SearchStreamService) sendPart(searchId string, data []byte) {
	s.hub.BroadcastMessage("SEARCH_RESULTS_"+searchId, data)
}

func (s *SearchStreamService) sendError(searchId string, err error) {
	payload := models.SearchErrorPayload{
		SearchId: searchId,
		Message:  err.Error(),
	}

	msg := ws_hub.WSMessage{
		Action: "SEARCH_ERROR",
		Data:   payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("json marshal error", logger.ErrAttr(err))
	}

	s.hub.BroadcastMessage("SEARCH_RESULTS_"+searchId, data)
}
