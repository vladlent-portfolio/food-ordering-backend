package user

import (
	"errors"
	"food_ordering_backend/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

const SessionCookieName = "access_token"
const ContextUserKey = "user"

type API struct {
	service *Service
}

func ProvideAPI(s *Service) *API {
	return &API{s}
}

func (api *API) Register(router *gin.RouterGroup, db *gorm.DB) {
	auth := InitAuthMiddleware(db)

	router.GET("", auth(true), api.FindAll)
	router.GET("/me", auth(false), api.Info)
	router.GET("/logout", auth(false), api.Logout)
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

	cookie := SessionCookie(session.Token, 0)

	http.SetCookie(c.Writer, cookie)
	c.Status(http.StatusOK)
}

func (api *API) Logout(c *gin.Context) {
	user := c.MustGet(ContextUserKey).(User)
	err := api.service.Logout(user)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	cookie := SessionCookie("", -1)

	http.SetCookie(c.Writer, cookie)
	c.Status(http.StatusOK)
}

func (api *API) Info(c *gin.Context) {
	user := c.MustGet(ContextUserKey).(User)
	c.JSON(http.StatusOK, ToResponseDTO(user))
}

func (api *API) bindAuthDTO(c *gin.Context) (AuthDTO, error) {
	var dto AuthDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return dto, err
	}
	return dto, nil
}

func SessionCookie(token string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   maxAge,
	}
}
