package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(SessionCookieName)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, err := service.FindSessionByToken(cookie.Value)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(ContextUserKey, session.User)
		c.Next()
	}
}
