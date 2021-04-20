package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthMid struct {
	userService *Service
	jwtService  *JWTService
}

func (m *AuthMid) Use(c *gin.Context) {
	fmt.Println("hello")
}

func ProvideAuthMid(service *Service, jwtService *JWTService) *AuthMid {
	return &AuthMid{service, jwtService}
}

func AuthMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(SessionCookieName)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwtService.AuthClaimsFromToken(cookie.Value)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(JWTUserIDKey, claims.UserID)
		c.Next()
	}
}
