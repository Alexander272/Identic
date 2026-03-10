package services

import (
	"context"
	"sort"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
)

type SearchService struct {
	repo repository.Search
}

func NewSearchService(repo repository.Search) *SearchService {
	return &SearchService{
		repo: repo,
	}
}

type Search interface {
	SearchAndGroup(ctx context.Context, req *models.SearchRequest) ([]*models.Results, error)
}

func (s *SearchService) SearchAndGroup(ctx context.Context, req *models.SearchRequest) ([]*models.Results, error) {
	if len(req.Items) == 0 {
		return nil, models.ErrNoData
	}

	for i := range req.Items {
		req.Items[i].Name = NormalizeString(req.Items[i].Name)
	}

	// 1. Получаем плоский список из БД
	rawResults := make([]*models.OrderMatchResult, 0)
	var err error
	if req.IsFuzzy {
		rawResults, err = s.repo.FindSimilar(ctx, req)
	} else {
		rawResults, err = s.repo.Find(ctx, req)
	}
	if err != nil {
		return nil, err
	}

	// 2. Группируем по годам в Go
	// Используем map для сбора, потом конвертируем в слайс
	yearMap := make(map[int][]*models.OrderMatchResult)
	for _, res := range rawResults {
		yearMap[res.Year] = append(yearMap[res.Year], res)
	}

	// 3. Формируем ответ
	var groups []*models.Results
	for year, apps := range yearMap {
		groups = append(groups, &models.Results{
			Year:   year,
			Orders: apps,
		})
	}

	// 4. Сортируем группы по году (например, от новых к старым)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Year > groups[j].Year
	})

	return groups, nil
}
