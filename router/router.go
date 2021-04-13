package router

import (
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	cat := r.Group("/categories")
	{
		catAPI := category.InitAPI(db)
		catAPI.Register(cat)
	}

	d := r.Group("/dishes")
	{
		dishesAPI := dish.InitAPI(db)
		dishesAPI.Register(d)
	}

	return r
}
