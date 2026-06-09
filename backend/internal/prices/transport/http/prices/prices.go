package prices

import (
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
	service services.Prices
}

func NewHandler(service services.Prices) *Handler {
	return &Handler{service: service}
}

func Register(api *gin.RouterGroup, service services.Prices, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	read := api.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourcePrice).Read()))
	read.GET("", handler.getAll)
	read.POST("search", handler.search)
	read.POST("search-all", handler.searchAll)

	write := api.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourcePrice).Write()))
	write.PUT("/batch", handler.batchSave)
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

	prices, total, err := h.service.GetAll(c.Request.Context(), page, perPage)
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

	prices, total, err := h.service.Search(c.Request.Context(), &req)
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

	prices, err := h.service.SearchAll(c.Request.Context(), &req)
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

	if err := h.service.BatchSave(c.Request.Context(), &req); err != nil {
		response.SendError(c, err, req)
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Цены сохранены"})
}
