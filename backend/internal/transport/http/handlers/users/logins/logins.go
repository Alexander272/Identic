package logins

import (
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.UserLogins
}

func NewHandler(service services.UserLogins) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.UserLogins, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	logins := api.Group("/logins", middleware.CheckPermissions(access.Reg.R(access.ResourceLogins).Read()))
	{
		logins.GET("/:id", handler.get)
	}
}

func (h *Handler) get(c *gin.Context) {
	strId := c.Param("id")
	id, err := uuid.Parse(strId)
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	req := &models.GetUserLoginsDTO{
		UserID: &id,
		Limit:  100,
	}

	data, err := h.service.GetByUser(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	response.SendData(c, data, len(data))
}
