package import_file

import (
	"net/http"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
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
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Не удалось получить файлы")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.NewErrorResponse(c, http.StatusNoContent, "no content", "Нет файлов для загрузки")
		return
	}

	sheetType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if files[0].Header.Get("Content-Type") != sheetType && !strings.Contains(files[0].Filename, "xls") {
		response.NewErrorResponse(c, http.StatusInternalServerError, "invalid type file", "Данный формат не поддерживается")
		return
	}

	dto := &models.ImportDTO{File: files[0]}
	if err := h.service.Load(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Файлы загружены") // logger.StringAttr("user_id", user.ID),
	// logger.StringAttr("username", user.Name),

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Файлы загружены"})
}
