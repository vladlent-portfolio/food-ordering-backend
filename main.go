package main

import (
	"food_ordering_backend/categories"
	"food_ordering_backend/db"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", pong)

	d := db.MustGet()
	d.AutoMigrate(&categories.Category{})

	r.Run()
}
func pong(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
