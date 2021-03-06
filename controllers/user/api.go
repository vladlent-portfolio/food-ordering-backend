package user

import (
	"errors"
	"food_ordering_backend/common"
	"food_ordering_backend/config"
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
	router.POST("/signup", api.Create)
	router.POST("/signin", api.Login)
}

// Create godoc
// @Summary Create new user
// @ID user-create
// @Tags user
// @Accept json
// @Param auth body AuthDTO true "User info"
// @Produce json
// @Success 201 {object} ResponseDTO
// @Failure 409,422
// @Router /users/signup [post]
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

	_, err = api.authorize(c, dto)

	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, ToResponseDTO(user))
}

// FindAll godoc
// @Summary Get all users. Requires admin rights.
// @ID user-get-all
// @Tags user
// @Param page query integer false "0-based page number"
// @Param limit query integer false "amount of entries per page"
// @Produce json
// @Success 200 {object} DTOsWithPagination
// @Failure 401
// @Router /users [get]
func (api *API) FindAll(c *gin.Context) {
	p := common.ExtractPagination(c, 10)
	users := api.service.FindAll(p)
	c.JSON(http.StatusOK, DTOsWithPagination{
		Users: ToResponseDTOs(users),
		Pagination: common.PaginationDTO{
			Page:  p.Page(),
			Limit: p.Limit(),
			Total: api.service.CountAll(),
		},
	})
}

// Login godoc
// @Summary Sign in
// @ID user-login
// @Tags user
// @Accept json
// @Param auth body AuthDTO true "User login data"
// @Success 200 {object} ResponseDTO
// @Failure 404,500
// @Router /users/signin [post]
func (api *API) Login(c *gin.Context) {
	dto, err := api.bindAuthDTO(c)

	if err != nil {
		return
	}

	session, err := api.authorize(c, dto)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, ToResponseDTO(session.User))
}

// Logout godoc
// @Summary Logout. Requires auth.
// @ID user-logout
// @Tags user
// @Success 200
// @Failure 401,500
// @Router /users/logout [get]
func (api *API) Logout(c *gin.Context) {
	// Skipping error because cookie was already checked in auth middleware.
	cookie, _ := c.Request.Cookie(SessionCookieName)

	err := api.service.Logout(cookie.Value)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, SessionCookie("", -1))
	c.Status(http.StatusOK)
}

// Info godoc
// @Summary Get info about current user. Requires auth.
// @ID user-info
// @Tags user
// @Success 200 {object} ResponseDTO
// @Failure 401
// @Router /users/me [get]
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

func (api *API) authorize(c *gin.Context, dto AuthDTO) (Session, error) {
	session, err := api.service.Login(dto)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrInvalidPassword):
			c.Status(http.StatusNotFound)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return Session{}, err
	}

	cookie := SessionCookie(session.Token, 0)
	http.SetCookie(c.Writer, cookie)

	return session, nil
}

func SessionCookie(token string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Domain:   config.ClientURL.Hostname(),
		Path:     "/",
		MaxAge:   maxAge,
	}
}
