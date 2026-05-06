package repository

import (
	"context"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/internal/repository/redis"
)

type Search interface {
	postgres.Search
	redis.Search
}

type searchProvider struct {
	redisRepo    redis.Search
	postgresRepo postgres.Search
}

func (s searchProvider) GetCache(ctx context.Context, req *models.GetCacheDTO) ([]string, error) {
	return s.redisRepo.GetCache(ctx, req)
}

func (s searchProvider) SetCache(ctx context.Context, req *models.SetCacheDTO) error {
	return s.redisRepo.SetCache(ctx, req)
}

func (s searchProvider) FetchExact(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	return s.postgresRepo.FetchExact(ctx, req)
}

func (s searchProvider) FetchFuzzy(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	return s.postgresRepo.FetchFuzzy(ctx, req)
}

func (s searchProvider) FetchExactByQuantity(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	return s.postgresRepo.FetchExactByQuantity(ctx, req)
}

func (s searchProvider) FetchFuzzyByQuantity(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	return s.postgresRepo.FetchFuzzyByQuantity(ctx, req)
}
