package category

import (
	"food_ordering_backend/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type API struct {
	Service *Service
}

func ProvideAPI(s *Service) *API {
	return &API{s}
}

func (api *API) Create(c *gin.Context) {
	var dto DTO
	err := c.BindJSON(&dto)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	category, err := api.Service.Create(ToCategory(dto))

	if err != nil {
		if common.IsDuplicateKeyErr(err) {
			c.Status(http.StatusConflict)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, ToDTO(category))
}

func (api *API) FindByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	cat, err := api.Service.FindByID(uint(id))

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, ToDTO(cat))
}

func (api *API) FindAll(c *gin.Context) {
	categories := api.Service.FindAll()
	c.JSON(http.StatusOK, ToDTOs(categories))
}
