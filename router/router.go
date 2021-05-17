package router

import (
	"food_ordering_backend/config"
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/order"
	"food_ordering_backend/controllers/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

type Controller interface {
	Register(g *gin.RouterGroup, db *gorm.DB)
}

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware())
	r.Static("/"+config.StaticDir, config.StaticDirAbs)

	routes := map[string]Controller{
		"/categories": category.InitAPI(db),
		"/dishes":     dish.InitAPI(db),
		"/users":      user.InitAPI(db),
		"/orders":     order.InitAPI(db),
	}

	for route, api := range routes {
		api.Register(r.Group(route), db)
	}

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedHeaders := []string{
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Cache-Control",
			"Connection",
			"Content-Type",
			"Content-Length",
		}

		c.Header("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
		c.Header("Access-Control", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
