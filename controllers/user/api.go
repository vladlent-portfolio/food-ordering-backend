package user

import (
	"errors"
	"food_ordering_backend/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

const SessionCookieName = "access_token"

type API struct {
	service *Service
}

func ProvideAPI(s *Service) *API {
	return &API{s}
}

func (api *API) Register(router *gin.RouterGroup) {
	router.GET("", api.FindAll)
	router.GET("/me", AuthMiddleware(), api.Info)
	router.POST("", api.Create)
	router.POST("/signin", api.Login)
}

func (api *API) Create(c *gin.Context) {
	dto, err := api.bindAuthDTO(c)

	if err != nil {
		return
	}

	user, err := api.service.Create(dto)

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
	users := api.service.FindAll()
	c.JSON(http.StatusOK, ToResponseDTOs(users))
}

func (api *API) Login(c *gin.Context) {
	dto, err := api.bindAuthDTO(c)

	if err != nil {
		return
	}

	session, err := api.service.Login(dto)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrInvalidPassword):
			c.Status(http.StatusForbidden)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    session.Token,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(c.Writer, cookie)
	c.Status(http.StatusOK)
}

func (api *API) Info(c *gin.Context) {
	uid := c.MustGet(JWTUserIDKey).(uint)
	u, err := api.service.FindByID(uid)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ToResponseDTO(u))
}

func (api *API) bindAuthDTO(c *gin.Context) (AuthDTO, error) {
	var dto AuthDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return dto, err
	}
	return dto, nil
}
