package import_file

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/prices/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Import
}

func NewHandler(service services.Import) *Handler {
	return &Handler{service: service}
}

func Register(api *gin.RouterGroup, service services.Import, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	write := api.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourcePrice).Write()))
	write.POST("/import", handler.importXLSX)
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

	if err := h.service.ImportXLSX(c.Request.Context(), io.NopCloser(bytes.NewReader(data))); err != nil {
		response.SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Цены импортированы"})
}
