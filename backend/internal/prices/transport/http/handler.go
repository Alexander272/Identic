package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/prices/services"
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
	prices := api.Group("v1/prices", middleware.VerifyToken)
	{
		prices.Use(middleware.CheckPermissions(access.Reg.R(access.ResourcePrice).Read()))
		prices.POST("search", h.search)
		prices.POST("/export", h.exportXLSX)

		prices.Use(middleware.CheckPermissions(access.Reg.R(access.ResourcePrice).Write()))
		prices.POST("/batch", h.batchSave)
		prices.POST("/import", h.importXLSX)
	}

}

func (h *Handler) search(c *gin.Context) {
	var req models.SearchPriceRequest
	if err := c.BindJSON(&req); err != nil {
		response.SendError(c, err)
		return
	}

	prices, total, err := h.services.Prices.Search(c.Request.Context(), req)
	if err != nil {
		response.SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse[[]*models.Price]{
		Data:  prices,
		Total: total,
	})
}

func (h *Handler) batchSave(c *gin.Context) {
	var req models.BatchSavePricesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, err)
		return
	}

	if err := h.services.Prices.BatchSave(c.Request.Context(), req); err != nil {
		response.SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Цены сохранены"})
}

func (h *Handler) exportXLSX(c *gin.Context) {
	var req models.ExportPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, err)
		return
	}

	data, err := h.services.Export.ExportXLSX(c.Request.Context(), req)
	if err != nil {
		response.SendError(c, err)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=prices.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *Handler) importXLSX(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.SendError(c, err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.SendError(c, err)
		return
	}

	if err := h.services.Import.ImportXLSX(c.Request.Context(), io.NopCloser(bytes.NewReader(data))); err != nil {
		response.SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Цены импортированы"})
}
