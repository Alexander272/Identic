package orders

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Orders
}

func NewHandler(service services.Orders) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Orders) {
	handler := NewHandler(service)

	orders := api.Group("/orders")
	{
		// orders.GET("", handler.getAll)
		orders.GET("/:id", handler.getById)
		orders.GET("/info/:id", handler.getInfo)
		orders.GET("/by-year/:year", handler.getByYear)
		orders.GET("/unique/:field", handler.getUniqueData)
		orders.GET("/flat", handler.getFlatData)
		orders.POST("", handler.create)
		orders.PUT("/:id", handler.update)
		// orders.DELETE("/:id", handler.delete)
	}
}

func (h *Handler) getById(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}
	req := &models.GetOrderByIdDTO{Id: id}

	order, err := h.service.GetById(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, &response.DataResponse{Data: order})
}

func (h *Handler) getInfo(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}
	req := &models.GetOrderByIdDTO{Id: id}

	positionIds := c.Query("positions")
	if positionIds != "" {
		req.PositionIds = strings.Split(positionIds, ",")
	}

	order, err := h.service.GetInfoById(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, &response.DataResponse{Data: order})
}

func (h *Handler) getByYear(c *gin.Context) {
	yearStr := c.Param("year")
	if yearStr == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "year is empty", "Отправлены некорректные данные")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	req := &models.GetOrderByYearDTO{Year: year}

	data, err := h.service.GetByYear(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) getUniqueData(c *gin.Context) {
	field := c.Param("field")
	if field == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "field is empty", "Отправлены некорректные данные")
		return
	}

	sort := c.Query("sort")
	if sort == "DESC" {
		sort = "DESC"
	} else {
		sort = "ASC"
	}

	req := &models.GetUniqueDTO{Field: field, Sort: sort}

	data, err := h.service.GetUniqueData(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) getFlatData(c *gin.Context) {
	req := &models.GetFlatOrderDTO{}

	for fields, val := range c.QueryMap("search") {
		req.Search = &models.Search{
			Value:  val,
			Fields: strings.Split(fields, ","),
		}
		break
	}

	sortLine := c.Query("sort")
	if sortLine != "" {
		sort, found := strings.CutPrefix(sortLine, "-")
		order := "ASC"
		if found {
			order = "DESC"
		}
		req.Sort = &models.Sort{
			Field: sort,
			Type:  order,
		}
	}

	cursorLine := c.Query("cursor")
	if cursorLine != "" && cursorLine != "null" {
		req.Cursor = cursorLine
	}

	limit := 50
	limitStr := c.Query("limit")
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	req.Page = &models.Page{
		Limit:  limit,
		Offset: 0,
	}

	data, err := h.service.GetFlatData(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.OrderDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, &response.IdResponse{Id: dto.Id, Message: "Заказ успешно создан"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}
	dto := &models.OrderDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}
	if id != dto.Id {
		response.NewErrorResponse(c, http.StatusBadRequest, "id is not equal to dto.Id", "Некорректные данные")
		return
	}

	if err := h.service.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, &response.IdResponse{Id: dto.Id, Message: "Заказ успешно обновлен"})
}
