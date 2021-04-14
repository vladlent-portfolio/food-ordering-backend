package main

import (
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"log"
)

func main() {
	db := database.MustGet()
	r := router.Setup(db)

	log.Panicln(r.Run())
}
