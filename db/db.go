package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var database *gorm.DB

// Init initializes db session.
// Successive calls to Init will do nothing if a successful connection has already been established.
func Init() error {
	if database != nil {
		return nil
	}

	db, err := gorm.Open(postgres.Open("postgres://vlad:123@localhost:5432/food-ordering-local?sslmode=disable"))

	if err != nil {
		return err
	}

	database = db
	return nil
}

// Get returns a db connection. Will call Init if necessary.
func Get() (*gorm.DB, error) {
	if database != nil {
		return database, nil
	}

	if err := Init(); err != nil {
		return database, nil
	}

	return database, nil
}

// MustGet same as Get but will panic if call to Init errors out.
func MustGet() *gorm.DB {
	db, err := Get()

	if err != nil {
		log.Panicln(err)
	}

	return db
}
