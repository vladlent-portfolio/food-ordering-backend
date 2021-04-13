package router

import (
	"food_ordering_backend/controllers/category"
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

	return r
}
