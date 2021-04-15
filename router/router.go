package router

import (
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type controller interface {
	Register(g *gin.RouterGroup)
}

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	routes := map[string]controller{
		"/categories": category.InitAPI(db),
		"/dishes":     dish.InitAPI(db),
	}

	for route, api := range routes {
		api.Register(r.Group(route))
	}

	return r
}
