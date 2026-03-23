package handlers

import (
	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/http/handlers/import_file"
	"github.com/Alexander272/Identic/backend/internal/transport/http/handlers/orders"
	"github.com/Alexander272/Identic/backend/internal/transport/http/handlers/search"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Services
	conf     *config.Config
	hub      *ws_hub.Hub
	// middleware *middleware.Middleware
}

type Deps struct {
	Services *services.Services
	Conf     *config.Config
	Hub      *ws_hub.Hub
	// Middleware *middleware.Middleware
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		services: deps.Services,
		conf:     deps.Conf,
		hub:      deps.Hub,
		// middleware: deps.Middleware,
	}
}

func (h *Handler) Init(group *gin.RouterGroup) {
	v1 := group.Group("/v1")

	// auth.Register(v1, auth.Deps{Service: h.services.Session, Middleware: h.middleware, Auth: h.conf.Auth})
	import_file.Register(v1, h.services.Import)
	search.Register(v1, h.services.SearchStream, h.hub)

	orders.Register(v1, h.services.Orders)

	// secure := v1.Group("", h.middleware.VerifyToken)
}
