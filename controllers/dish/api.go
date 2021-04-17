package dish

import (
	"food_ordering_backend/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type API struct {
	service *Service
}

func ProvideAPI(s *Service) *API {
	return &API{s}
}

func (api *API) Register(router *gin.RouterGroup) {
	router.GET("", api.FindAll)
	router.GET("/:id", api.FindByID)
	router.POST("", api.Create)
	router.PUT("/:id", api.Update)
	router.DELETE("/:id", api.Delete)
}

func (api *API) Create(c *gin.Context) {
	dto, err := api.bindJSON(c)

	if err != nil {
		return
	}

	dish, err := api.service.Create(ToModel(dto))

	if err != nil {
		if common.IsDuplicateKeyErr(err) {
			c.Status(http.StatusConflict)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, ToDTO(dish))
}

func (api *API) FindByID(c *gin.Context) {
	dish, err := api.findByID(c)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, ToDTO(dish))
}

func (api *API) FindAll(c *gin.Context) {
	var dishes []Dish
	cidq := c.Query("cid")

	if cidq == "" {
		dishes = api.service.FindAll(0)
	} else {
		cid, err := strconv.ParseUint(cidq, 10, 64)

		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		dishes = api.service.FindAll(uint(cid))
	}

	c.JSON(http.StatusOK, ToDTOs(dishes))
}

func (api *API) Update(c *gin.Context) {
	dish, err := api.findByID(c)

	if err != nil {
		return
	}

	dto, err := api.bindJSON(c)

	if err != nil {
		return
	}

	dish.Title = dto.Title
	dish.Price = dto.Price
	dish.CategoryID = dto.CategoryID
	dish.Category.ID = dto.CategoryID

	dish, err = api.service.Save(dish)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ToDTO(dish))
}

func (api *API) Delete(c *gin.Context) {
	dish, err := api.findByID(c)

	if err != nil {
		return
	}

	dish, err = api.service.Delete(dish)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ToDTO(dish))
}

func (api *API) bindJSON(c *gin.Context) (DTO, error) {
	var dto DTO
	err := c.BindJSON(&dto)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return dto, err
	}

	return dto, nil
}

func (api *API) findByID(c *gin.Context) (Dish, error) {
	var dish Dish
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return dish, err
	}

	dish, err = api.service.FindByID(uint(id))

	if err != nil {
		c.Status(http.StatusNotFound)
		return dish, err
	}

	return dish, nil
}
