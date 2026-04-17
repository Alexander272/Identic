package search

import (
	"context"
	"net/http"
	"time"

	"github.com/Alexander272/Identic/backend/internal/constants"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.SearchStream
	hub     *ws_hub.Hub
}

func NewHandler(service services.SearchStream, hub *ws_hub.Hub) *Handler {
	return &Handler{
		service: service,
		hub:     hub,
	}
}

func Register(api *gin.RouterGroup, service services.SearchStream, hub *ws_hub.Hub) {
	handler := NewHandler(service, hub)

	search := api.Group("/search")
	{
		search.POST("", handler.search)
		search.POST("/stream", handler.search)
	}
}

func (h *Handler) search(c *gin.Context) {
	dto := &models.SearchRequest{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}
	dto.SearchId = uuid.NewString()

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	dto.ActorID = user.ID
	dto.ActorName = user.Name

	go func() {
		ctx := context.Background()
		topic := "SEARCH_RESULTS_" + dto.SearchId

		ready := h.hub.WaitForFirstSubscriber(ctx, topic, 10*time.Second)

		if ready {
			h.service.Streaming(c, dto)
		} else {
			logger.Info("Search cancelled: no subscribers found", logger.StringAttr("search_id", dto.SearchId))
		}
	}()

	c.JSON(http.StatusOK, &response.IdResponse{Id: dto.SearchId, Message: "Запрос успешно отправлен"})
}
