package order

import (
	"errors"
	"fmt"
	"food_ordering_backend/controllers/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type API struct {
	service *Service
}

func ProvideAPI(s *Service) *API {
	return &API{s}
}

func (api *API) Register(router *gin.RouterGroup, db *gorm.DB) {
	auth := user.InitAuthMiddleware(db)

	router.GET("", auth(false), api.FindAll)
	router.POST("", auth(false), api.Create)
	router.PATCH("/:id/cancel", auth(false), api.Cancel)
	router.PUT("/:id", auth(true), api.Update)
}

func (api *API) FindAll(c *gin.Context) {
	var orders []Order
	var err error
	u := c.MustGet(user.ContextUserKey).(user.User)

	if u.IsAdmin {
		orders, err = api.service.FindAll()
	} else {
		orders, err = api.service.FindByUID(u.ID)
	}

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ToResponseDTOs(orders))
}

func (api *API) Create(c *gin.Context) {
	var dto RequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	fmt.Printf("dto: %+v\n", dto)

	u := c.MustGet(user.ContextUserKey).(user.User)
	o, err := api.service.Create(dto.Items, u)

	if err != nil {
		var errDishID *ErrDishID
		if errors.As(err, &errDishID) {
			c.String(http.StatusBadRequest, errDishID.Error())
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, ToResponseDTO(o))
}

func (api *API) Cancel(c *gin.Context) {

}

func (api *API) Update(c *gin.Context) {}
