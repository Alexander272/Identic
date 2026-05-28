package middleware

import (
	"fmt"

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
			response.SendError(c, models.ErrSessionEmpty)
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
			response.SendError(c, fmt.Errorf("%w: %v", models.ErrPolicyCheck, lastErr))
			return
		}

		if !accessAllowed {
			response.SendError(c, models.ErrPermissionDenied)
			return
		}

		c.Next()
	}
}
