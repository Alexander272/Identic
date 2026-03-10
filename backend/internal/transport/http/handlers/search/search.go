package search

import (
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Search
}

func NewHandler(service services.Search) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Search) {
	handler := NewHandler(service)

	search := api.Group("/search")
	{
		search.POST("", handler.search)
	}
}

func (h *Handler) search(c *gin.Context) {
	dto := &models.SearchRequest{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	data, err := h.service.SearchAndGroup(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, &response.DataResponse{Data: data})
}
