package postgres

var Tables = struct {
	Prices          string
	PriceSearchLogs string
}{
	Prices:          "prices",
	PriceSearchLogs: "price_search_logs",
}
