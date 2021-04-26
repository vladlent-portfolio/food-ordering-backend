// +build wireinject

package order

import (
	"food_ordering_backend/controllers/dish"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitAPI(db *gorm.DB) *API {
	wire.Build(ProvideAPI, ProvideService, ProvideRepository, dish.ServiceSet)
	return nil
}
