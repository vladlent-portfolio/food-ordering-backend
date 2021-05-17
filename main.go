package main

import (
	"fmt"
	"food_ordering_backend/config"
	"food_ordering_backend/database"
	_ "food_ordering_backend/docs"
	"food_ordering_backend/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Food Ordering Backend
// @version 1.0
// @description Golang backend for Food Ordering portfolio app.
// @contact.name Vladlen Tereshchenko
// @contact.url https://github.com/VladlenT
// @contact.email vladlent.dev@gmail.com
// @license.name MIT

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	db := database.MustGet()
	r := router.Setup(db)

	url := ginSwagger.URL(config.HostRaw + "/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	var address = fmt.Sprintf("%s:%s", config.HostURL.Hostname(), config.HostURL.Port())
	log.Panicln(r.Run(address))
}
