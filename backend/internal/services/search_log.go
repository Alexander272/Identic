package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/Alexander272/Identic/backend/pkg/logger"
)

type SearchLogService struct {
	repo repository.SearchLogs
}

func NewSearchLogService(repo repository.SearchLogs) *SearchLogService {
	return &SearchLogService{repo: repo}
}

type SearchLogRecorder interface {
	Get(ctx context.Context, dto *models.GetSearchLogsDTO) ([]*models.SearchLog, error)
	Create(ctx context.Context, dto *models.CreateSearchLogDTO) error
	LogAsync(req *models.SearchRequest, originalItems []models.SearchItem, duration time.Duration, resultsCount int)
}

func (s *SearchLogService) Create(ctx context.Context, dto *models.CreateSearchLogDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create search log: %w", err)
	}
	return nil
}

func (s *SearchLogService) Get(ctx context.Context, dto *models.GetSearchLogsDTO) ([]*models.SearchLog, error) {
	data, err := s.repo.Get(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get search logs: %w", err)
	}
	return data, nil
}

func (s *SearchLogService) LogAsync(req *models.SearchRequest, originalItems []models.SearchItem, duration time.Duration, resultsCount int) {
	var searchType models.SearchType
	switch {
	case req.SearchByQuantityOnly && req.IsFuzzy:
		searchType = models.SearchTypeQuantityFuzzy
	case req.SearchByQuantityOnly:
		searchType = models.SearchTypeQuantityExact
	case req.IsFuzzy:
		searchType = models.SearchTypeFuzzy
	default:
		searchType = models.SearchTypeExact
	}

	dto := &models.CreateSearchLogDTO{
		SearchId:     req.SearchId,
		ActorID:      req.ActorID,
		ActorName:    req.ActorName,
		SearchType:   searchType,
		Query:        originalItems,
		DurationMs:   duration.Milliseconds(),
		ResultsCount: resultsCount,
		ItemsCount:   len(originalItems),
	}

	go func() {
		if err := s.Create(context.Background(), dto); err != nil {
			logger.Error("failed to log search", logger.ErrAttr(err))
			error_bot.Send(nil, fmt.Sprintf("failed to log search. error: %v", err), req)
		}
	}()
}
