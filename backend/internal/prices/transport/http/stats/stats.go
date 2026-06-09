package stats

import (
	"strconv"
	"time"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/prices/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.PriceSearchLogs
}

func NewHandler(service services.PriceSearchLogs) *Handler {
	return &Handler{service: service}
}

func Register(api *gin.RouterGroup, service services.PriceSearchLogs, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	read := api.Group("statistics", middleware.CheckPermissions(access.Reg.R(access.ResourcePrice).Read()))
	read.GET("search", handler.getSearchLogs)
}

func (h *Handler) getSearchLogs(c *gin.Context) {
	var dto models.GetPriceSearchLogsDTO

	if actorId := c.Query("actorId"); actorId != "" {
		id, err := uuid.Parse(actorId)
		if err != nil {
			response.SendError(c, err)
			return
		}
		dto.ActorID = &id
	}
	if startDate := c.Query("startDate"); startDate != "" {
		t, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			response.SendError(c, err)
			return
		}
		dto.StartDate = &t
	}
	if endDate := c.Query("endDate"); endDate != "" {
		t, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			response.SendError(c, err)
			return
		}
		dto.EndDate = &t
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			response.SendError(c, err)
			return
		}
		dto.Limit = limit
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			response.SendError(c, err)
			return
		}
		dto.Offset = offset
	}

	logs, err := h.service.Get(c.Request.Context(), &dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}

	response.SendData(c, logs, len(logs))
}
