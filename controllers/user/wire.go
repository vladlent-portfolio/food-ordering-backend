// +build wireinject

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var set = wire.NewSet(ProvideService, ProvideRepository, ProvideJWTService)

func InitAPI(db *gorm.DB) *API {
	wire.Build(ProvideAPI, set)
	return nil
}

func ProvideAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	wire.Build(AuthMiddleware, set)
	return nil
}
