package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/go-redis/redis/v8"
)

type SearchRepo struct {
	db *redis.Client
}

func NewSearchRepo(db *redis.Client) *SearchRepo {
	return &SearchRepo{
		db: db,
	}
}

type Search interface {
	GetCache(ctx context.Context, req *models.GetCacheDTO) ([]string, error)
	SetCache(ctx context.Context, req *models.SetCacheDTO) error
}

func (r *SearchRepo) GetCache(ctx context.Context, req *models.GetCacheDTO) ([]string, error) {
	str, err := r.db.Get(ctx, fmt.Sprintf("%s_%s", req.SearchId, req.OrderId)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, models.ErrSearchExpired
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return strings.Split(str, ";"), nil
}

func (r *SearchRepo) SetCache(ctx context.Context, req *models.SetCacheDTO) error {
	key := fmt.Sprintf("%s_%s", req.SearchId, req.OrderId)
	err := r.db.Set(ctx, key, strings.Join(req.PositionIds, ";"), req.Exp).Err()
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
