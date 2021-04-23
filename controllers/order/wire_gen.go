// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package order

import (
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitAPI(db *gorm.DB) *API {
	repository := ProvideRepository(db)
	service := ProvideService(repository)
	api := ProvideAPI(service)
	return api
}
