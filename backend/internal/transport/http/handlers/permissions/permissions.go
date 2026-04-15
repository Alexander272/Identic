package permissions

import (
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Permissions
}

func NewHandler(service services.Permissions) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Permissions, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	permissions := api.Group("/permissions", middleware.CheckPermissions(access.Reg.R(access.ResourcePerm).Read()))
	{
		permissions.GET("resources", handler.getResources)
	}
}

func (h *Handler) getResources(c *gin.Context) {
	data := h.service.GetResources(c)
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}
