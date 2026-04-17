package stats

import (
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Statistic
}

func NewHandler(service services.Statistic) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Statistic, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	stats := api.Group("/statistics")
	{
		search := stats.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceSearch).Read()))
		{
			search.GET("/search", handler.getSearch)
		}

		activity := stats.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceActivity).Read()))
		{
			activity.GET("/activity", handler.getActivity)
		}
	}
}

func (h *Handler) getSearch(c *gin.Context) {
	req := &models.GetSearchLogsDTO{}

	data, err := h.service.GetSearch(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) getActivity(c *gin.Context) {
	req := &models.GetAllActivityLogsDTO{}

	data, err := h.service.GetActivity(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}
