package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/prices/repository"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/Alexander272/Identic/backend/pkg/logger"
)

type PriceSearchLogs interface {
	Get(ctx context.Context, dto *models.GetPriceSearchLogsDTO) ([]*models.PriceSearchLog, error)
	Create(ctx context.Context, dto *models.CreatePriceSearchLogDTO) error
	LogAsync(codes, queries []string, actorID uuid.UUID, actorName string, duration time.Duration, resultsCount int)
}

type PriceSearchLogService struct {
	repo repository.PriceSearchLogs
}

func NewPriceSearchLogService(repo repository.PriceSearchLogs) *PriceSearchLogService {
	return &PriceSearchLogService{repo: repo}
}

func (s *PriceSearchLogService) Create(ctx context.Context, dto *models.CreatePriceSearchLogDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create price search log: %w", err)
	}
	return nil
}

func (s *PriceSearchLogService) Get(ctx context.Context, dto *models.GetPriceSearchLogsDTO) ([]*models.PriceSearchLog, error) {
	data, err := s.repo.Get(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get price search logs: %w", err)
	}
	return data, nil
}

func (s *PriceSearchLogService) LogAsync(codes, queries []string, actorID uuid.UUID, actorName string, duration time.Duration, resultsCount int) {
	dto := &models.CreatePriceSearchLogDTO{
		Queries:      queries,
		Codes:        codes,
		ActorID:      actorID,
		ActorName:    actorName,
		DurationMs:   duration.Milliseconds(),
		ResultsCount: resultsCount,
	}

	go func() {
		if err := s.Create(context.Background(), dto); err != nil {
			logger.Error("failed to log price search", logger.ErrAttr(err))
			error_bot.Send(nil, fmt.Sprintf("failed to log price search. error: %v", err), nil)
		}
	}()
}
