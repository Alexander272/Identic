package models

// SearchItem - позиция, которую мы ищем (входящий запрос)
type SearchItem struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
}

// SearchRequest - весь запрос на поиск похожего заказа
type SearchRequest struct {
	Items   []SearchItem `json:"items"`
	IsFuzzy bool         `json:"isFuzzy"`
}

// MatchedPosition - структура, соответствующая таблице positions
type MatchedPosition struct {
	Id         string  `json:"id" db:"id"`
	OrderId    string  `json:"orderId" db:"order_id"`
	Name       string  `json:"name" db:"name"`
	Quantity   float64 `json:"quantity" db:"quantity"`
	InputName  string  `json:"inputName" db:"input_name"`
	InputQty   float64 `json:"inputQty" db:"input_qty"`
	Similarity float64 `json:"similarity" db:"similarity"` // Заполняется SQL запросом
}

// OrderMatchResult - результат поиска по одному заказу
type OrderMatchResult struct {
	OrderId      string  `json:"orderId"`
	Customer     string  `json:"customer"`
	Consumer     string  `json:"consumer"`
	Year         int     `json:"year"`
	Link         string  `json:"link"`
	Score        float64 `json:"score"`        // Общий процент совпадения (0-100)
	MatchedCount int     `json:"matchedCount"` // Сколько позиций совпало
	TotalCount   int     `json:"totalCount"`   // Сколько позиций в запросе
}

type Results struct {
	Year   int                 `json:"year"`
	Orders []*OrderMatchResult `json:"orders"`
}
