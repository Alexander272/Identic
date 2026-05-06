package orders

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/constants"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
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

func Register(api *gin.RouterGroup, service services.Orders, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	orders := api.Group("/orders", middleware.CheckPermissions(access.Reg.R(access.ResourceOrder).Read()))
	{
		orders.GET("", handler.get)
		orders.GET("/:id", handler.getById)
		orders.GET("/info/:id", handler.getInfo)
		orders.GET("/by-year/:year", handler.getByYear)
		orders.GET("/unique/:field", handler.getUniqueData)
		orders.GET("/flat", handler.getFlatData)

		write := orders.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceOrder).Write()))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
		}

		// orders.DELETE("/:id", handler.delete)
	}
}

func (h *Handler) get(c *gin.Context) {
	req := &models.OrderFilterDTO{}

	filters := c.QueryMap("filters")

	for k, v := range filters {
		valueMap := c.QueryMap(k)
		values := []*models.FilterValue{}
		for key, value := range valueMap {
			values = append(values, &models.FilterValue{
				CompareType: key,
				Value:       value,
			})
		}
		if values[0].CompareType == "" {
			continue
		}

		f := &models.Filter{
			Field:     k,
			FieldType: v,
			Values:    values,
		}

		req.Filters = append(req.Filters, f)
	}

	data, err := h.service.Get(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, &response.DataResponse{Data: data})
}

func (h *Handler) getById(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}
	req := &models.GetOrderByIdDTO{Id: id}

	search := c.Query("search")
	if search != "" {
		req.SearchId = search
	}

	order, err := h.service.GetById(c, nil, req)
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

	// positionIds := c.Query("positions")
	// if positionIds != "" {
	// 	req.PositionIds = strings.Split(positionIds, ",")
	// }
	search := c.Query("search")
	if search != "" {
		req.SearchId = search
	}

	order, err := h.service.GetInfoById(c, req)
	if err != nil {
		if errors.Is(err, models.ErrNoData) {
			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), "")
			return
		}
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

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	dto.Actor = models.Actor{
		ID:   user.ID,
		Name: user.Name,
	}

	id, err := h.service.Create(c, dto)
	if err != nil {
		if errors.Is(err, models.ErrOrderAlreadyExists) {
			// response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Заказ уже существует")
			c.JSON(http.StatusConflict, &response.IdResponse{Id: id, Message: "Заказ уже существует"})
			return
		}

		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	// if id != "" {
	// 	c.JSON(http.StatusConflict, &response.IdResponse{Id: id, Message: "Заказ уже существует"})
	// 	return
	// }

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

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	dto.Actor = models.Actor{
		ID:   user.ID,
		Name: user.Name,
	}

	if err := h.service.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, &response.IdResponse{Id: dto.Id, Message: "Заказ успешно обновлен"})
}
