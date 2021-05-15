package category

import (
	"fmt"
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"food_ordering_backend/controllers/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type API struct {
	service *Service
}

func ProvideAPI(s *Service) *API {
	return &API{s}
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
// @Summary Create new category. Requires admin rights.
// @ID category-create
// @Tags category
// @Accept json
// @Param dto body DTO true "Category DTO"
// @Produce json
// @Success 201 {object} DTO
// @Failure 409,422,500
// @Router /categories [post]
func (api *API) Create(c *gin.Context) {
	dto, err := api.bindJSON(c)

	if err != nil {
		return
	}

	category, err := api.service.Create(ToModel(dto))

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

// FindByID godoc
// @Summary Find category by id
// @ID category-find
// @Tags category
// @Param id path integer true "Category id"
// @Produce json
// @Success 200 {object} DTO
// @Failure 403,404
// @Router /categories/:id [get]
func (api *API) FindByID(c *gin.Context) {
	cat, err := api.findByID(c)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, ToDTO(cat))
}

// FindAll godoc
// @Summary Get all categories
// @ID category-all
// @Tags category
// @Produce json
// @Success 200 {array} DTO
// @Failure 403,404
// @Router /categories [get]
func (api *API) FindAll(c *gin.Context) {
	categories := api.service.FindAll()
	c.JSON(http.StatusOK, ToDTOs(categories))
}

// Update godoc
// @Summary Replace category. Requires admin rights.
// @ID category-update
// @Tags category
// @Accept json
// @Param dto body DTO true "Category DTO"
// @Param id path integer true "Category id"
// @Produce json
// @Success 200 {object} DTO
// @Failure 401,403,404,500
// @Router /categories/:id [put]
func (api *API) Update(c *gin.Context) {
	cat, err := api.findByID(c)

	if err != nil {
		return
	}

	dto, err := api.bindJSON(c)

	if err != nil {
		return
	}

	cat.Title = dto.Title
	cat.Removable = dto.Removable

	cat, err = api.service.Save(cat)

	if err != nil {
		if common.IsDuplicateKeyErr(err) {
			c.Status(http.StatusConflict)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, ToDTO(cat))
}

// Upload godoc
// @Summary Upload image for category. Requires admin rights.
// @ID category-upload
// @Tags category
// @Param id path integer true "Category id"
// @Param image formData file true "Category image"
// @Accept multipart/form-data
// @Produce text/plain
// @Success 200 {string} string "Link to uploaded image"
// @Failure 400,401,404,413,415,500
// @Router /categories/:id/upload [patch]
func (api *API) Upload(c *gin.Context) {
	cat, err := api.findByID(c)

	if err != nil {
		return
	}

	fileHeader, err := c.FormFile("image")

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if fileHeader.Size > config.MaxUploadFileSize {
		c.Status(http.StatusRequestEntityTooLarge)
		return
	}

	file, err := fileHeader.Open()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	contentType, err := common.MIMEType(file)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	if contentType != "image/png" && contentType != "image/jpeg" {
		c.Status(http.StatusUnsupportedMediaType)
		return
	}

	fPath := filepath.Join(config.CategoriesImgDirAbs, fmt.Sprintf("%d.%s", cat.ID, path.Base(contentType)))
	if err := os.MkdirAll(filepath.Dir(fPath), os.ModeDir); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := c.SaveUploadedFile(fileHeader, fPath); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	fName := filepath.Base(fPath)
	cat.Image = &fName

	cat, err = api.service.Save(cat)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, PathToImg(fName))
}

// Delete godoc
// @Summary Delete category by id. Requires admin rights.
// @ID category-delete
// @Tags category
// @Param id path integer true "Category id"
// @Produce json
// @Success 200 {object} DTO
// @Failure 401,403,404,500
// @Router /categories/:id [delete]
func (api *API) Delete(c *gin.Context) {
	cat, err := api.findByID(c)

	if err != nil {
		return
	}

	if !cat.Removable {
		c.Status(http.StatusForbidden)
		return
	}

	cat, err = api.service.Delete(cat)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ToDTO(cat))
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

func (api *API) findByID(c *gin.Context) (Category, error) {
	var cat Category
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return cat, err
	}

	cat, err = api.service.FindByID(uint(id))

	if err != nil {
		c.Status(http.StatusNotFound)
		return cat, err
	}

	return cat, nil
}
