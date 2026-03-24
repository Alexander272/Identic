package services

import (
	"context"
	"runtime"

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
	results, err := s.service.Search(ctx, req)
	if err != nil {
		s.sendError(req.SearchId, err)
	}

	const batchSize = 10
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

		msg := ws_hub.WSMessage{
			Action: "SEARCH_RESULT_PART",
			Data:   payload,
		}

		data, err := json.Marshal(msg)
		if err != nil {
			s.sendError(req.SearchId, err)
		}

		s.sendPart(req.SearchId, data)

		runtime.Gosched()
	}
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
