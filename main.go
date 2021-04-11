package main

import (
	"food_ordering_backend/category"
	"food_ordering_backend/database"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.MustGet()
	db.AutoMigrate(&category.Category{})

	categoryAPI := category.InitAPI(db)

	r := gin.Default()

	r.GET("/categories", categoryAPI.FindAll)
	r.POST("/categories", categoryAPI.Create)

	r.Run()
}
