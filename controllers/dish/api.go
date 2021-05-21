package dish

import (
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"food_ordering_backend/controllers/user"
	"food_ordering_backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type API struct {
	service *Service
	upload  *services.Upload
}

func ProvideAPI(s *Service) *API {
	upload := &services.Upload{
		AllowedTypes: []string{"image/png", "image/jpeg"},
		MaxFileSize:  config.MaxUploadFileSize,
		Root:         config.DishesImgDirAbs,
		FormDataKey:  "image",
	}
	return &API{s, upload}
}

func (api *API) Register(router *gin.RouterGroup, db *gorm.DB) {
	auth := user.InitAuthMiddleware(db)

	router.GET("", api.FindAll)
	router.GET("/:id", api.FindByID)
	router.POST("", auth(true), api.Create)
	router.PUT("/:id", auth(true), api.Update)
	router.PATCH("/:id/upload", auth(true), api.Upload)
	router.DELETE("/:id", auth(true), api.Delete)
}

// Create godoc
// @Summary Create new dish. Requires admin rights.
// @ID dish-create
// @Tags dish
// @Accept json
// @Param dto body DTO true "User info"
// @Produce json
// @Success 201 {object} DTO
// @Failure 409,422
// @Router /dishes [post]
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

// FindByID godoc
// @Summary Find dish by id
// @ID dish-find
// @Tags dish
// @Param id path integer true "Dish id"
// @Produce json
// @Success 200 {object} DTO
// @Failure 403,404
// @Router /dishes/:id [get]
func (api *API) FindByID(c *gin.Context) {
	dish, err := api.findByID(c)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, ToDTO(dish))
}

// FindAll godoc
// @Summary Get all dishes
// @ID dish-all
// @Tags dish
// @Param cid query integer false "filter dishes by category id"
// @Produce json
// @Success 200 {array} DTO
// @Failure 403,404
// @Router /dishes [get]
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

// Update godoc
// @Summary Replace dish. Requires admin rights.
// @ID dish-update
// @Tags dish
// @Accept json
// @Param dto body DTO true "Dish DTO"
// @Param id path integer true "Dish id"
// @Produce json
// @Success 200 {object} DTO
// @Failure 401,403,404,409,500
// @Router /dishes/:id [put]
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
		if common.IsDuplicateKeyErr(err) {
			c.Status(http.StatusConflict)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, ToDTO(dish))
}

// Upload godoc
// @Summary Upload image for dish. Requires admin rights.
// @ID dish-upload
// @Tags dish
// @Param id path integer true "Dish id"
// @Param image formData file true "Dish image"
// @Accept multipart/form-data
// @Produce text/plain
// @Success 200 {string} string "Link to uploaded image"
// @Failure 400,401,404,413,415,500
// @Router /dishes/:id/upload [patch]
func (api *API) Upload(c *gin.Context) {
	dish, err := api.findByID(c)

	if err != nil {
		return
	}

	absPath := api.upload.ParseAndSave(c, strconv.Itoa(int(dish.ID)))

	if absPath == "" {
		return
	}

	name := filepath.Base(absPath)
	dish.Image = &name

	dish, err = api.service.Save(dish)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, PathToImg(name))
}

// Delete godoc
// @Summary Delete dish by id. Requires admin rights.
// @ID dish-delete
// @Tags dish
// @Param id path integer true "Dish id"
// @Produce json
// @Success 200 {object} DTO
// @Failure 401,403,404,500
// @Router /dishes/:id [delete]
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

	dto.Title = strings.TrimSpace(dto.Title)

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
