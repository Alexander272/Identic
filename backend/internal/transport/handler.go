package http

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/http/handlers"
	"github.com/Alexander272/Identic/backend/internal/transport/ws"
	"github.com/Alexander272/Identic/backend/pkg/auth"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/Alexander272/Identic/backend/pkg/limiter"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	keycloak *auth.KeycloakClient
	services *services.Services
	hub      *ws_hub.Hub
}

func NewHandler(keycloak *auth.KeycloakClient, services *services.Services, hub *ws_hub.Hub) *Handler {
	return &Handler{
		keycloak: keycloak,
		services: services,
		hub:      hub,
	}
}

func (h *Handler) Init(conf *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		limiter.Limit(conf.Limiter.RPS, conf.Limiter.Burst, conf.Limiter.TTL),
		gin.CustomRecovery(h.ErrorHandler),
	)

	router.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router, conf)

	return router
}

func (h *Handler) ErrorHandler(c *gin.Context, origErr any) {
	err := fmt.Errorf("unexpected error: %v", origErr)
	error_bot.Send(c, err.Error(), gin.H{"PANIC": true, "Stack trace": string(debug.Stack())})
	debug.PrintStack()
	response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла непредвиденная ошибка: "+err.Error())
}

func (h *Handler) initAPI(router *gin.Engine, conf *config.Config) {
	// middleware := middleware.NewMiddleware(h.services, conf.Auth, h.keycloak)
	// handler := handlers.NewHandler(&handlers.Deps{Services: h.services, Conf: conf, Middleware: middleware})
	handler := handlers.NewHandler(&handlers.Deps{Services: h.services, Conf: conf, Hub: h.hub})
	wsHandler := ws.NewWsHandler(h.hub, conf.Http, h.services)

	api := router.Group("/api")
	handler.Init(api)

	api.GET("/ws", func(c *gin.Context) {
		wsHandler.HandleWS(c)
	})
}
