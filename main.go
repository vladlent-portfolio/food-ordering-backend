package main

import (
	"food_ordering_backend/category"
	"food_ordering_backend/database"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.MustGet()
	db.AutoMigrate(&category.Category{})

	categoryAPI := &category.API{&category.Service{&category.Repository{DB: db}}}

	r := gin.Default()

	r.POST("/categories", categoryAPI.Create)
	r.GET("/categories", categoryAPI.FindAll)

	r.Run()
}
