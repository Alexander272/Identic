package users

import (
	"net/http"

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
	service services.Users
}

func NewHandler(service services.Users) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Users, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	users := api.Group("/users", middleware.CheckPermissions(access.Reg.R(access.ResourceUser).Read()))
	{
		users.GET("", handler.getAll)
		// 	read.GET("/access", handler.getByAccess)
		// 	read.GET("/realm/:id", handler.getByRealm)
		// 	read.GET("/:id", handler.getById)
		// 	read.GET("/sso/:id", handler.getBySSOId)

		write := users.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceUser).Write()))
		{
			write.POST("/sync", handler.sync)
			// write.POST("", handler.create)
			// write.POST("/several", handler.createSeveral)
			write.PUT("/:id", handler.update)
		}
	}
}

func (h *Handler) getAll(c *gin.Context) {
	data, err := h.service.GetAll(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

// func (h *Handler) getById(c *gin.Context) {
// 	id := c.Param("id")
// 	err := uuid.Validate(id)
// 	if err != nil {
// 		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
// 		return
// 	}

// 	data, err := h.service.GetById(c, id)
// 	if err != nil {
// 		if errors.Is(err, models.ErrNoRows) {
// 			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), err.Error())
// 			return
// 		}
// 		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
// 		error_bot.Send(c, err.Error(), id)
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.DataResponse{Data: data})
// }

// func (h *Handler) getBySSOId(c *gin.Context) {
// 	id := c.Param("id")
// 	err := uuid.Validate(id)
// 	if err != nil {
// 		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
// 		return
// 	}

// 	data, err := h.service.GetBySSOId(c, id)
// 	if err != nil {
// 		if errors.Is(err, models.ErrNoRows) {
// 			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), err.Error())
// 			return
// 		}
// 		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
// 		error_bot.Send(c, err.Error(), id)
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.DataResponse{Data: data})
// }

func (h *Handler) sync(c *gin.Context) {
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	actor := &models.Actor{
		ID:   user.ID,
		Name: user.Name,
	}

	if err := h.service.Sync(c, actor); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Пользователи синхронизированы"})
}

// func (h *Handler) create(c *gin.Context) {
// 	dto := &models.UserData{}
// 	if err := c.BindJSON(dto); err != nil {
// 		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
// 		return
// 	}

// 	if err := h.service.Create(c, dto); err != nil {
// 		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
// 		error_bot.Send(c, err.Error(), dto)
// 		return
// 	}
// 	c.JSON(http.StatusCreated, response.IdResponse{Message: "Пользователь создан"})
// }

// func (h *Handler) createSeveral(c *gin.Context) {
// 	var dto []*models.UserData
// 	if err := c.BindJSON(&dto); err != nil {
// 		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
// 		return
// 	}

// 	if err := h.service.CreateSeveral(c, nil, dto); err != nil {
// 		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
// 		error_bot.Send(c, err.Error(), dto)
// 		return
// 	}
// 	c.JSON(http.StatusCreated, response.IdResponse{Message: "Пользователи созданы"})
// }

func (h *Handler) update(c *gin.Context) {
	dto := &models.UserDataDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	strId := c.Param("id")
	id, err := uuid.Parse(strId)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
		return
	}
	if id != dto.ID {
		response.NewErrorResponse(c, http.StatusBadRequest, "id is not equal to dto.ID", "Некорректные данные")
		return
	}
	dto.ID = id

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
	c.JSON(http.StatusOK, response.IdResponse{Message: "Пользователь обновлен"})
}
