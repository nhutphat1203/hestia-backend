package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nhutphat1203/hestia-backend/internal/domain"
	"github.com/nhutphat1203/hestia-backend/pkg/errorf"
	"github.com/nhutphat1203/hestia-backend/pkg/response"
)

func AuthMiddleware(authenticator domain.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.SendError(c, errorf.HttpStatus(errorf.Unauthorized), errorf.Message(errorf.Unauthorized), errorf.Unauthorized)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")

		ok, err := authenticator.Authenticate(token)
		if !ok || err != nil {
			response.SendError(c, errorf.HttpStatus(errorf.InvalidToken), errorf.Message(errorf.InvalidToken), errorf.InvalidToken)
			return
		}

		c.Next()
	}
}
