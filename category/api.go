package category

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	category := api.Service.Create(ToCategory(dto))

	c.JSON(http.StatusOK, ToDTO(category))
}

func (api *API) FindAll(c *gin.Context) {
	categories := api.Service.FindAll()

	c.JSON(http.StatusOK, ToCategoriesDTOs(categories))
}
