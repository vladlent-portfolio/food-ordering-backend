package user

import (
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(SessionCookieName)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})

		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*AuthClaims); ok && claims.Valid(jwt.DefaultValidationHelper) == nil {
			c.Set(JWTUserIDKey, claims.UserID)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
