// +build wireinject

package dish

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

var ServiceSet = wire.NewSet(ProvideService, ProvideRepository)

func InitAPI(db *gorm.DB) *API {
	wire.Build(ProvideAPI, ServiceSet)
	return nil
}
