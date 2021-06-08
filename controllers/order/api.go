package order

import (
	"errors"
	"food_ordering_backend/common"
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
	router.PATCH("/:id", auth(true), api.Patch)
	router.PUT("/:id", auth(true), api.Update)
}

// FindAll godoc
// @Summary Get all orders. Requires auth.
// @Description If requester is admin, it returns all orders. Otherwise, it returns orders only for that user.
// @ID order-all
// @Tags order
// @Param page query integer false "0-based page number"
// @Param limit query integer false "amount of entries per page"
// @Produce json
// @Success 200 {object} DTOsWithPagination
// @Failure 401,403,404,500
// @Router /orders [get]
func (api *API) FindAll(c *gin.Context) {
	var orders []Order
	var err error
	u := c.MustGet(user.ContextUserKey).(user.User)
	p := common.ExtractPagination(c, 10)

	if u.IsAdmin {
		u.ID = 0
	}

	orders, err = api.service.FindAll(u.ID, p)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, DTOsWithPagination{
		Orders: ToResponseDTOs(orders),
		Pagination: common.PaginationDTO{
			Page:  p.Page(),
			Limit: p.Limit(),
			Total: api.service.CountAll(u.ID),
		},
	})
}

// Create godoc
// @Summary Create new order. Requires auth.
// @ID order-create
// @Tags order
// @Accept json
// @Param dto body CreateDTO true "Create order DTO"
// @Produce json
// @Success 201 {object} ResponseDTO
// @Failure 401,403,422,500
// @Router /orders [post]
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

// Patch godoc
// @Summary Patch order. Requires admin rights.
// @ID order-patch
// @Tags order
// @Param status query integer true "New order status"
// @Success 204
// @Failure 401,403,404,500
// @Router /orders/:id [patch]
func (api *API) Patch(c *gin.Context) {
	s := c.Query("status")
	status, err := strconv.Atoi(s)

	if err != nil || !IsValidStatus(status) {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	o, err := api.findByID(c)

	if err != nil {
		return
	}

	if o.Status == Status(status) {
		c.Status(http.StatusNotModified)
		return
	}

	if err := api.service.UpdateStatus(o.ID, Status(status)); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

// Update godoc
// @Summary Replace order. Requires admin rights.
// @ID order-update
// @Tags order
// @Accept json
// @Param dto body UpdateDTO true "Order update DTO"
// @Param id path integer true "Order id"
// @Produce json
// @Success 200 {object} ResponseDTO
// @Failure 401,403,404,500
// @Router /orders/:id [put]
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
