package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthMiddlewareFunc func(isAdmin bool) gin.HandlerFunc

func ProvideAuthMiddleware(service *Service) AuthMiddlewareFunc {
	return func(isAdmin bool) gin.HandlerFunc {
		return AuthMiddleware(service, isAdmin)
	}
}

// AuthMiddleware intercepts a request to check whether user is authorized
// and has a valid session token. If any of checks fails request will be
// automatically aborted with http.StatusUnauthorized.
// If check is successful, middleware will set ContextUserKey with
// authorized User in current gin.Context.
// adminOnly flag indicates that a user must have admin rights to access the route.
func AuthMiddleware(service *Service, adminOnly bool) gin.HandlerFunc {
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

		if adminOnly && !session.User.IsAdmin {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(ContextUserKey, session.User)
		c.Next()
	}
}
