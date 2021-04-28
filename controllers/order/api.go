package order

import (
	"errors"
	"food_ordering_backend/controllers/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	var dto CreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
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
	o, err := api.findByID(c)

	if err != nil {
		return
	}

	u := c.MustGet(user.ContextUserKey).(user.User)

	if !u.IsAdmin && u.ID != o.UserID {
		c.Status(http.StatusForbidden)
		return
	}

	switch o.Status {
	case StatusCanceled, StatusDone:
		c.Status(http.StatusNotModified)
		return
	}

	if err := api.service.UpdateStatus(o.ID, StatusCanceled); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (api *API) Update(c *gin.Context) {
	o, err := api.findByID(c)

	if err != nil {
		return
	}

	var dto UpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	o, err = api.service.Update(o, dto)

	if err != nil {
		var errDishID *ErrDishID
		if errors.As(err, &errDishID) {
			c.String(http.StatusBadRequest, errDishID.Error())
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ToResponseDTO(o))
}

func (api *API) findByID(c *gin.Context) (Order, error) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return Order{}, err
	}

	o, err := api.service.FindByID(uint(id))

	if err != nil {
		var errOrderID *ErrOrderID

		if errors.As(err, &errOrderID) {
			c.String(http.StatusNotFound, errOrderID.Error())
		} else {
			c.Status(http.StatusInternalServerError)
		}

		return Order{}, err
	}

	return o, nil
}
