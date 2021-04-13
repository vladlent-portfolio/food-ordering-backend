// +build wireinject

package category

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitAPI(db *gorm.DB) *API {
	wire.Build(ProvideAPI, ProvideService, ProvideRepository)
	return nil
}
