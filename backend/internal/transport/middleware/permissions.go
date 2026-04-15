package middleware

import (
	"net/http"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/constants"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/gin-gonic/gin"
)

type Permission struct {
	Section string
	Method  string
}

func (m *Middleware) CheckPermissions(required ...access.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, exists := c.Get(constants.CtxUser)
		if !exists {
			response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "сессия не найдена")
			c.Abort()
			return
		}
		user := u.(models.User)

		var accessAllowed bool
		var lastErr error

		for _, r := range required {
			ok, err := m.services.AccessPolices.Enforce(user.ID.String(), string(r.Resource), string(r.Action))
			if err != nil {
				lastErr = err
				continue
			}
			if ok {
				accessAllowed = true
				break
			}
		}

		if lastErr != nil && !accessAllowed {
			response.NewErrorResponse(c, http.StatusInternalServerError, lastErr.Error(), "ошибка проверки прав")
			c.Abort()
			return
		}

		if !accessAllowed {
			response.NewErrorResponse(c, http.StatusForbidden, "forbidden", "недостаточно прав")
			c.Abort()
			return
		}

		c.Next()
	}
}
