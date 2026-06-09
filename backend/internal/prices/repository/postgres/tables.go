package postgres

var Tables = struct {
	Prices          string
	PriceSearchLogs string

	Users string
}{
	Prices:          "prices",
	PriceSearchLogs: "price_search_logs",

	Users: "users",
}
