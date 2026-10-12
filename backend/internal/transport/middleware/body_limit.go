package middleware

import (
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/gin-gonic/gin"
)

// MaxBodyBytes — максимальный размер тела запроса (16 МБ).
// Покрывает импорт больших XLSX и поиск по десяткам тысяч кодов,
// но отсекает аномально большие тела
const MaxBodyBytes int64 = 16 << 20

// BodyLimit ограничивает размер тела запроса: сначала проверяет Content-Length
// (чтобы вернуть понятную ошибку до чтения тела), затем ограничивает само чтение
// через http.MaxBytesReader — на случай запросов без Content-Length (chunked)
func (m *Middleware) BodyLimit(c *gin.Context) {
	if c.Request.ContentLength > MaxBodyBytes {
		response.SendError(c, models.ErrRequestTooLarge)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	c.Next()
}
