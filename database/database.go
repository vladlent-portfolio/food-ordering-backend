package database

import (
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/order"
	"food_ordering_backend/controllers/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var database *gorm.DB

// Init initializes db session.
// Successive calls to Init will do nothing if a successful connection has already been established.
func Init() error {
	connStr := "postgres://vlad:123@localhost:5432/food-ordering-local?sslmode=disable"
	return initDB(connStr)
}

func InitTest() error {
	connStr := "postgres://vlad:123@localhost:5432/food-ordering-test?sslmode=disable"
	return initDB(connStr)
}

// Get returns a db connection. Will call Init if necessary.
func Get() (*gorm.DB, error) {
	return get(false)
}

func GetTest() (*gorm.DB, error) {
	return get(true)
}

// MustGet same as Get but will panic if call to Init errors out.
func MustGet() *gorm.DB {
	db, err := Get()

	if err != nil {
		log.Panicln(err)
	}

	return db
}

func MustGetTest() *gorm.DB {
	db, err := GetTest()

	if err != nil {
		log.Panicln(err)
	}

	return db
}

func initDB(connStr string) error {
	if database != nil {
		return nil
	}

	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return err
	}

	autoMigrate(db)

	database = db
	return nil
}

func autoMigrate(db *gorm.DB) {
	models := []interface{}{
		&category.Category{},
		&dish.Dish{},
		&user.User{},
		&user.Session{},
		&order.Order{},
		&order.Item{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			panic(err)
		}
	}

}

func get(isTest bool) (*gorm.DB, error) {
	if database != nil {
		return database, nil
	}
	var err error

	if isTest {
		err = InitTest()
	} else {
		err = Init()
	}

	if err != nil {
		return nil, err
	}

	return database, nil
}
