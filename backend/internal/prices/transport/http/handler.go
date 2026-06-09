package http

import (
	"github.com/Alexander272/Identic/backend/internal/prices/services"
	"github.com/Alexander272/Identic/backend/internal/prices/transport/http/export"
	"github.com/Alexander272/Identic/backend/internal/prices/transport/http/import_file"
	"github.com/Alexander272/Identic/backend/internal/prices/transport/http/prices"
	"github.com/Alexander272/Identic/backend/internal/prices/transport/http/stats"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init(api *gin.RouterGroup, middleware *middleware.Middleware) {
	pricesGroup := api.Group("v1/prices", middleware.VerifyToken)

	prices.Register(pricesGroup, h.services.Prices, middleware)
	export.Register(pricesGroup, h.services.Export, middleware)
	import_file.Register(pricesGroup, h.services.Import, middleware)
	stats.Register(pricesGroup, h.services.PriceSearchLogs, middleware)
}
