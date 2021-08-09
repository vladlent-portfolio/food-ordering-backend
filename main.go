package main

import (
	_ "embed"
	"fmt"
	"food_ordering_backend/config"
	"food_ordering_backend/database"
	"food_ordering_backend/docs"
	"food_ordering_backend/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

//go:embed static/index.html
var indexHTML []byte

//go:embed static/robots.txt
var robotsTxt []byte

// @title Food Ordering Backend
// @version 1.0
// @description Golang backend for Food Ordering app.
// @contact.name Vladlen Tereshchenko
// @contact.url https://github.com/VladlenT
// @contact.email vladlent.dev@gmail.com
// @license.name MIT

// @BasePath /

func main() {
	updateSwaggerDoc()

	db := database.MustGet()
	r := router.Setup(db)

	r.GET("/", serveEmbedded("text/html", indexHTML))
	r.GET("/robots.txt", serveEmbedded("text/plain", robotsTxt))

	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	log.Panicln(r.Run(":" + viper.GetString("HOST_PORT")))
}

func updateSwaggerDoc() {
	hostURL := config.HostURL

	if config.IsProdMode {
		docs.SwaggerInfo.Schemes = []string{"https"}
		docs.SwaggerInfo.Host = hostURL.Hostname()

	} else {
		docs.SwaggerInfo.Schemes = []string{"http"}
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", hostURL.Hostname(), hostURL.Port())
	}

	docs.SwaggerInfo.Description = fmt.Sprintf(`Golang backend for Food Ordering App.
Frontend available [here](%s).

*Since Swagger doesn't support cookie-based authorizations you should **Sign In** in the [frontend app](%s) and then comeback here to be able to interact with guarded routes.* `,
		config.ClientURL.String(), config.ClientURL.String(),
	)
}

func serveEmbedded(contentType string, data []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(200, contentType, data)
	}
}
