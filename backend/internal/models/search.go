package models

import (
	"time"

	"github.com/google/uuid"
)

// SearchItem - позиция, которую мы ищем (входящий запрос)
type SearchItem struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
}

// SearchRequest - весь запрос на поиск похожего заказа
type SearchRequest struct {
	Items     []SearchItem `json:"items"`
	IsFuzzy   bool         `json:"isFuzzy"`
	ActorID   uuid.UUID    `json:"actorId"`
	ActorName string       `json:"actorName"`
	SearchId  string
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
	OrderId      string    `json:"orderId"`
	Customer     string    `json:"customer"`
	Consumer     string    `json:"consumer"`
	Date         time.Time `json:"date"`
	Year         int       `json:"year"`
	Link         string    `json:"link"`
	Score        float64   `json:"score"`        // Общий процент совпадения (0-100)
	MatchedPos   int       `json:"matchedPos"`   // Сколько позиций совпало
	MatchedQuant int       `json:"matchedQuant"` // Сколько позиций + количество совпало
	TotalCount   int       `json:"totalCount"`   // Сколько позиций в запросе
	// PositionIds  []string `json:"positionIds"`
	Positions []*MatchPosition `json:"positions"`
}
type MatchPosition struct {
	Id         string `json:"id"`
	ReqId      string `json:"reqId"`
	QuantEqual bool   `json:"quantEqual"`
}

type Results struct {
	Year   int                 `json:"year"`
	Count  int                 `json:"count"`
	Orders []*OrderMatchResult `json:"orders"`
}

type MatchInfo struct {
	PosID     string
	ReqID     string
	ItemScore float64 // Совокупный балл (текст + количество)
	ReqQty    float64
	DbQty     float64
}

type OrderInfo struct {
	Id       string
	Year     int
	Customer string
	Consumer string
	Matches  map[string]MatchInfo
}

type RawMatch struct {
	OrderId    string // UUID из базы
	YearInt    int
	Customer   string
	Consumer   string
	Date       time.Time
	ReqId      string   // ID позиции из запроса (номер по порядку)
	PosId      string   // UUID позиции в базе
	PSearch    string   // Название товара в базе
	ReqTokens  []string // Токены (только для Fuzzy)
	ReqQty     float64  // Сколько просил пользователь
	DbQty      float64  // Сколько реально в базе
	Similarity float64  // 1.0 для точного поиска, 0.3-1.0 для Fuzzy
}

type SearchResultPart struct {
	Items  []*OrderMatchResult `json:"items"`
	IsLast bool                `json:"isLast"`
	Total  int                 `json:"total"`
}

type SearchErrorPayload struct {
	SearchId string `json:"searchId"`
	Message  string `json:"message"`
}

type GetCacheDTO struct {
	OrderId  string `json:"orderId"`
	SearchId string `json:"searchId"`
}

type SetCacheDTO struct {
	OrderId     string   `json:"orderId"`
	SearchId    string   `json:"searchId"`
	PositionIds []string `json:"positionIds"`
	Exp         time.Duration
}
