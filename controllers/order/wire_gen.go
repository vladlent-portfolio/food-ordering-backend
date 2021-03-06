// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package order

import (
	"food_ordering_backend/controllers/dish"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitAPI(db *gorm.DB) *API {
	repository := ProvideRepository(db)
	dishRepository := dish.ProvideRepository(db)
	service := dish.ProvideService(dishRepository)
	orderService := ProvideService(repository, service)
	api := ProvideAPI(orderService)
	return api
}
