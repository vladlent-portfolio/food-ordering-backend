package user

import (
	"errors"
	"food_ordering_backend/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

const SessionCookieName = "session-id"

type API struct {
	Service *Service
}

func ProvideAPI(s *Service) *API {
	return &API{s}
}

func (api *API) Register(router *gin.RouterGroup) {
	router.GET("", api.FindAll)
	router.POST("", api.Create)
	router.POST("/signin", api.Login)
}

func (api *API) Create(c *gin.Context) {
	var dto AuthDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	user, err := api.Service.Create(dto)

	if err != nil {
		switch {
		case common.IsDuplicateKeyErr(err):
			c.Status(http.StatusConflict)
		case errors.Is(err, gorm.ErrInvalidValue):
			c.String(http.StatusUnprocessableEntity, err.Error())
		default:
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, ToResponseDTO(user))
}

func (api *API) FindAll(c *gin.Context) {
	users := api.Service.FindAll()
	c.JSON(http.StatusOK, ToResponseDTOs(users))
}

func (api *API) Login(c *gin.Context) {

}
