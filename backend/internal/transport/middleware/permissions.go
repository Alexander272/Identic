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

		for _, r := range required {
			accessAllowed, err := m.services.AccessPolices.Enforce(user.ID, string(r.Resource), string(r.Action))

			if err != nil {
				response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
				c.Abort()
				return
			}

			if !accessAllowed {
				response.NewErrorResponse(c, http.StatusForbidden, "forbidden", "недостаточно прав")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
