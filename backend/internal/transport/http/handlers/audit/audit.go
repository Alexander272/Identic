package audit

import (
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.AuditLogs
}

func NewHandler(service services.AuditLogs) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.AuditLogs, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	audit := api.Group("/audit")
	{
		audit.GET("", handler.get)
	}
}

func (h *Handler) get(c *gin.Context) {
	data, err := h.service.Get(c, &models.GetAuditLogsDTO{})
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}
