package export

import (
	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/prices/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Export
}

func NewHandler(service services.Export) *Handler {
	return &Handler{service: service}
}

func Register(api *gin.RouterGroup, service services.Export, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	read := api.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourcePrice).Read()))
	read.POST("/export", handler.exportXLSX)
}

func (h *Handler) exportXLSX(c *gin.Context) {
	var req models.ExportPriceRequest
	if err := c.BindJSON(&req); err != nil {
		response.SendError(c, err)
		return
	}

	data, err := h.service.ExportXLSX(c.Request.Context(), req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}

	response.SendFile(c, data, "prices.xlsx")
}
