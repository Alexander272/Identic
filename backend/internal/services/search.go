package services

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
)

type SearchService struct {
	repo     repository.Search
	orderUrl string
	cacheTTL time.Duration
}

func NewSearchService(repo repository.Search, orderUrl string, cacheTTL time.Duration) *SearchService {
	return &SearchService{
		repo:     repo,
		orderUrl: orderUrl,
		cacheTTL: cacheTTL,
	}
}

type Search interface {
	Search(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error)
	SearchAndGroup(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error)
	GetCache(ctx context.Context, req *models.GetCacheDTO) ([]string, error)
}

func (s *SearchService) Search(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
	var rawMatches []*models.RawMatch
	var err error

	for i := range req.Items {
		req.Items[i].Name = NormalizeString(req.Items[i].Name)
	}

	// 1. Получаем данные из репозитория
	if req.IsFuzzy {
		rawMatches, err = s.repo.FetchFuzzy(ctx, req)
	} else {
		rawMatches, err = s.repo.FetchExact(ctx, req)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find orders. error: %w", err)
	}

	orderMap := make(map[string]*models.OrderMatchResult)
	bestMatches := make(map[string]map[string]*models.MatchInfo)

	for _, m := range rawMatches {
		if req.IsFuzzy && !s.validateTokens(m.PSearch, m.ReqTokens) {
			continue
		}

		// Расчет Score (для точного поиска similarity придет из БД как 1.0)
		itemScore := s.calculateItemScore(m.Similarity, m.ReqQty, m.DbQty)

		if _, ok := orderMap[m.OrderId]; !ok {
			orderMap[m.OrderId] = &models.OrderMatchResult{
				OrderId:  m.OrderId,
				Year:     m.YearInt,
				Customer: m.Customer,
				Consumer: m.Consumer,
				Date:     m.Date,
			}
			bestMatches[m.OrderId] = make(map[string]*models.MatchInfo)
		}

		// Сохраняем только лучшее совпадение (актуально для неточного поиска)
		if prev, ok := bestMatches[m.OrderId][m.ReqId]; !ok || itemScore > prev.ItemScore {
			bestMatches[m.OrderId][m.ReqId] = &models.MatchInfo{
				PosID:     m.PosId,
				ReqID:     m.ReqId,
				ItemScore: itemScore,
				ReqQty:    m.ReqQty,
				DbQty:     m.DbQty,
			}
		}
	}

	return s.finalize(orderMap, bestMatches, len(req.Items), req.SearchId), nil
}

func (s *SearchService) SearchAndGroup(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
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
		return nil, fmt.Errorf("failed to find orders. error: %w", err)
	}

	// генерируем ссылку на заказ
	for _, m := range rawResults {
		if err := s.genLink(m, req.SearchId); err != nil {
			return nil, err
		}
	}

	// 2. Группируем по годам в Go
	// Используем map для сбора, потом конвертируем в слайс
	// yearMap := make(map[int][]*models.OrderMatchResult)
	// for _, res := range rawResults {
	// 	yearMap[res.Year] = append(yearMap[res.Year], res)
	// }

	// // 3. Формируем ответ
	// var groups []*models.Results
	// for year, apps := range yearMap {
	// 	groups = append(groups, &models.Results{
	// 		Year:   year,
	// 		Orders: apps,
	// 		Count:  len(apps),
	// 	})
	// }

	// // 4. Сортируем группы по году (например, от новых к старым)
	// sort.Slice(groups, func(i, j int) bool {
	// 	return groups[i].Year > groups[j].Year
	// })

	return rawResults, nil
}

func (s *SearchService) GetCache(ctx context.Context, req *models.GetCacheDTO) ([]string, error) {
	positions, err := s.repo.GetCache(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache. error: %w", err)
	}
	return positions, nil
}

// Вспомогательные методы логики
func (s *SearchService) genLink(item *models.OrderMatchResult, searchId string) (err error) {
	item.Link, err = url.JoinPath(s.orderUrl, item.OrderId)
	if err != nil {
		return fmt.Errorf("failed to generate link. error: %w", err)
	}

	url, err := url.Parse(item.Link)
	if err != nil {
		return fmt.Errorf("failed to parse link. error: %w", err)
	}
	q := url.Query()
	ids := make([]string, len(item.Positions))
	for i, p := range item.Positions {
		ids[i] = p.Id
	}
	// q.Set("positions", strings.Join(ids, ","))
	// url.RawQuery = q.Encode()

	cacheDTO := &models.SetCacheDTO{OrderId: item.OrderId, SearchId: searchId, PositionIds: ids, Exp: s.cacheTTL}
	s.repo.SetCache(context.Background(), cacheDTO)

	q.Set("search", searchId)
	url.RawQuery = q.Encode()

	item.Link = url.String()

	return nil
}

func (s *SearchService) validateTokens(pSearch string, reqTokens []string) bool {
	isValid := true
	pSearchLower := strings.ToLower(pSearch)
	dbTokens := s.makeTokenMap(pSearchLower)

	matchedSpecificCount := 0
	totalSpecificCount := 0

	for _, token := range reqTokens {
		// Считаем токен "специфичным", если в нем есть цифры или он длинный
		isSpecific := s.containsDigits(token) || len(token) > 5

		if isSpecific {
			totalSpecificCount++
			if _, exists := dbTokens[token]; exists {
				matchedSpecificCount++
			} else {
				// Если специфичный токен (например, часть артикула) не найден - это мусор
				isValid = false
				break
			}
		} else if len(token) >= 2 {
			// Для коротких слов (типа "Г" или "В") требуем точного совпадения,
			// но не бракуем весь результат сразу, если это общее слово.
			// Но если это тип СНП (А, В, Г), он критичен.
			if s.isCriticalType(token) {
				if _, exists := dbTokens[token]; !exists {
					isValid = false
					break
				}
			}
		}
	}

	if !isValid || (totalSpecificCount > 0 && matchedSpecificCount == 0) {
		return false
	}
	return true
}

func (s *SearchService) makeTokenMap(searchStr string) map[string]struct{} {
	words := strings.Fields(strings.ToLower(searchStr))
	m := make(map[string]struct{}, len(words))
	for _, w := range words {
		m[w] = struct{}{}
	}
	return m
}
func (s *SearchService) containsDigits(str string) bool {
	return strings.ContainsAny(str, "0123456789")
}

var reSymbols = regexp.MustCompile(`[а-яА-Яa-zA-Z]`)

// Критичные токены - это одиночные буквы, которые часто значат тип (А, В, Г, П)
func (s *SearchService) isCriticalType(str string) bool {
	if len(str) != 1 {
		return false
	}
	// Проверяем, что это буква, а не просто мусор
	return reSymbols.MatchString(str)
}

func (s *SearchService) calculateItemScore(sml, reqQty, dbQty float64) float64 {
	ratio := math.Min(reqQty, dbQty) / math.Max(reqQty, dbQty)
	qtyFactor := 0.8 + (0.2 * ratio)

	return sml * qtyFactor
}

func (s *SearchService) finalize(
	om map[string]*models.OrderMatchResult, bm map[string]map[string]*models.MatchInfo, total int, searchId string,
) []*models.OrderMatchResult {
	results := make([]*models.OrderMatchResult, 0, len(om))
	for id, order := range om {
		var sumScore float64
		var fullMatches int
		for _, m := range bm[id] {
			sumScore += m.ItemScore
			// order.PositionIds = append(order.PositionIds, m.PosID)
			order.Positions = append(order.Positions, &models.MatchPosition{
				Id:         m.PosID,
				ReqId:      m.ReqID,
				QuantEqual: m.DbQty == m.ReqQty,
			})

			if math.Abs(m.ReqQty-m.DbQty) < 0.0001 {
				fullMatches++
			}
		}
		order.MatchedPos = len(bm[id])
		order.MatchedQuant = fullMatches
		order.TotalCount = total
		order.Score = math.Round((sumScore/float64(total)*100)*100) / 100

		s.genLink(order, searchId)

		if order.Score > 0 {
			results = append(results, order)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Date.After(results[j].Date)
	})

	return results
}
