// +build wireinject

package order

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitAPI(db *gorm.DB) *API {
	wire.Build(ProvideAPI, ProvideService, ProvideRepository)
	return nil
}
