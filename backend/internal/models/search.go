package models

// SearchItem - позиция, которую мы ищем (входящий запрос)
type SearchItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// SearchRequest - весь запрос на поиск похожего заказа
type SearchRequest struct {
	Items []SearchItem `json:"items"`
}

// DBPosition - структура, соответствующая таблице positions
type DBPosition struct {
	ID         int     `db:"id"`
	OrderID    int     `db:"order_id"`
	Name       string  `db:"name"`
	Quantity   int     `db:"quantity"`
	Similarity float64 `db:"similarity"` // Заполняется SQL запросом
}

// OrderMatchResult - результат поиска по одному заказу
type OrderMatchResult struct {
	OrderID      int     `json:"order_id"`
	Score        float64 `json:"score"`         // Общий процент совпадения (0-100)
	MatchedCount int     `json:"matched_count"` // Сколько позиций совпало
	TotalCount   int     `json:"total_count"`   // Сколько позиций в запросе
}
