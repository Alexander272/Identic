package import_file

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Import
}

func NewHandler(service services.Import) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Import) {
	handler := NewHandler(service)

	importFile := api.Group("/import")
	{
		importFile.POST("", handler.load)
	}
}

func (h *Handler) load(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.SendError(c, err)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, "files is empty"))
		return
	}

	sheetType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if files[0].Header.Get("Content-Type") != sheetType && !strings.Contains(files[0].Filename, "xls") {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, "invalid type file"))
		return
	}

	dto := &models.ImportDTO{File: files[0]}
	if err := h.service.Load(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Файлы загружены") // logger.StringAttr("user_id", user.ID),
	// logger.StringAttr("username", user.Name),

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Файлы загружены"})
}
