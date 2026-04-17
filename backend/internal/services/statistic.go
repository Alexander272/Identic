package services

import (
	"context"

	"github.com/Alexander272/Identic/backend/internal/models"
)

type StatisticService struct {
	activity Activity
	search   SearchLogRecorder
}

func NewStatisticService(activity Activity, search SearchLogRecorder) *StatisticService {
	return &StatisticService{
		activity: activity,
		search:   search,
	}
}

type Statistic interface {
	GetSearch(ctx context.Context, dto *models.GetSearchLogsDTO) ([]*models.SearchLog, error)
	GetActivity(ctx context.Context, dto *models.GetAllActivityLogsDTO) ([]*models.ActivityLog, error)
}

func (s *StatisticService) GetSearch(ctx context.Context, dto *models.GetSearchLogsDTO) ([]*models.SearchLog, error) {
	return s.search.Get(ctx, dto)
}

func (s *StatisticService) GetActivity(ctx context.Context, dto *models.GetAllActivityLogsDTO) ([]*models.ActivityLog, error) {
	return s.activity.Get(ctx, dto)
}
