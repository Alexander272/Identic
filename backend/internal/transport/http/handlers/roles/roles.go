package roles

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/constants"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Roles
}

func NewHandler(service services.Roles) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Roles, middleware *middleware.Middleware) {
	handlers := NewHandler(service)

	roles := api.Group("/roles", middleware.CheckPermissions(access.Reg.R(access.ResourceRole).Read()))
	{
		roles.GET("", handlers.getAll)
		roles.GET("/all/stats", handlers.getWithStats)

		write := roles.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceRole).Write()))
		{
			write.GET("/item/:name", handlers.get)

			write.POST("", handlers.create)
			write.PUT("/:id", handlers.update)
		}

		delete := roles.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceRole).Delete()))
		{
			delete.DELETE("/:id", handlers.delete)
		}

		permissions := roles.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourcePerm).Read()))
		{
			permissions.GET("/:id/permissions", handlers.getWithPermissions)
		}
	}
}

func (h *Handler) getAll(c *gin.Context) {
	roles, err := h.service.GetAll(c)
	if err != nil {
		response.SendError(c, err)
		return
	}

	response.SendData(c, roles, len(roles))
}

func (h *Handler) get(c *gin.Context) {
	slug := c.Param("name")
	dto := &models.GetRoleDTO{Slug: slug}

	role, err := h.service.GetOne(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}

	response.SendData(c, role)
}

func (h *Handler) getWithStats(c *gin.Context) {
	roles, err := h.service.GetWithStats(c)
	if err != nil {
		response.SendError(c, err)
		return
	}

	response.SendData(c, roles, len(roles))
}

func (h *Handler) getWithPermissions(c *gin.Context) {
	strId := c.Param("id")
	id, err := uuid.Parse(strId)
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	role, err := h.service.GetOneWithPermissions(c, &models.GetRoleDTO{ID: id})
	if err != nil {
		response.SendError(c, err, id)
		return
	}
	response.SendData(c, role)
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.RoleDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, err)
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	dto.Actor = models.Actor{
		ID:   user.ID,
		Name: user.Name,
	}

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Роль создана"})
}

func (h *Handler) update(c *gin.Context) {
	strId := c.Param("id")
	id, err := uuid.Parse(strId)
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.RoleDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, err)
		return
	}
	if id != dto.ID {
		response.SendError(c, fmt.Errorf("%w: %s", models.ErrInvalidInput, "id is not equal to dto.ID"))
		return
	}
	dto.ID = id

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	dto.Actor = models.Actor{
		ID:   user.ID,
		Name: user.Name,
	}

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Роль обновлена"})
}

func (h *Handler) delete(c *gin.Context) {
	strId := c.Param("id")
	id, err := uuid.Parse(strId)
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	dto := &models.DeleteRoleDTO{ID: id}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	dto.Actor = models.Actor{
		ID:   user.ID,
		Name: user.Name,
	}

	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	c.Status(http.StatusNoContent)
}
