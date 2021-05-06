package router

import (
	"food_ordering_backend/config"
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/order"
	"food_ordering_backend/controllers/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller interface {
	Register(g *gin.RouterGroup, db *gorm.DB)
}

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// TODO: Fix before production
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Next()
	})

	r.Static("/static", config.StaticDirAbs)

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
