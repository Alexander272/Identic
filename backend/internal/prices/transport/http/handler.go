package http

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/constants"
	base_models "github.com/Alexander272/Identic/backend/internal/models"
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
		prices.GET("", h.getAll)
		prices.POST("search", h.search)
		prices.POST("search-all", h.searchAll)
		prices.POST("/export", h.exportXLSX)

		prices.Use(middleware.CheckPermissions(access.Reg.R(access.ResourcePrice).Write()))
		prices.POST("/batch", h.batchSave)
		prices.POST("/import", h.importXLSX)
	}

}

func (h *Handler) getAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		perPage = 20
	}

	prices, total, err := h.services.Prices.GetAll(c.Request.Context(), page, perPage)
	if err != nil {
		response.SendError(c, err, gin.H{"page": page, "limit": perPage})
		return
	}

	response.SendData(c, prices, total)
}

func (h *Handler) search(c *gin.Context) {
	var req models.SearchPriceRequest
	if err := c.BindJSON(&req); err != nil {
		response.SendError(c, err)
		return
	}

	if user, ok := c.Get(constants.CtxUser); ok {
		if u, ok := user.(base_models.User); ok {
			req.ActorID = u.ID
			req.ActorName = u.Name
		}
	}

	prices, total, err := h.services.Prices.Search(c.Request.Context(), &req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}

	response.SendData(c, prices, total)
}

func (h *Handler) searchAll(c *gin.Context) {
	var req models.SearchPriceRequest
	if err := c.BindJSON(&req); err != nil {
		response.SendError(c, err)
		return
	}

	prices, err := h.services.Prices.SearchAll(c.Request.Context(), &req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}

	response.SendData(c, prices, len(prices))
}

func (h *Handler) batchSave(c *gin.Context) {
	var req models.BatchSavePricesRequest
	if err := c.BindJSON(&req); err != nil {
		response.SendError(c, err)
		return
	}

	if err := h.services.Prices.BatchSave(c.Request.Context(), &req); err != nil {
		response.SendError(c, err, req)
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Цены сохранены"})
}

func (h *Handler) exportXLSX(c *gin.Context) {
	var req models.ExportPriceRequest
	if err := c.BindJSON(&req); err != nil {
		response.SendError(c, err)
		return
	}

	data, err := h.services.Export.ExportXLSX(c.Request.Context(), req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}

	response.SendFile(c, data, "prices.xlsx")
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
