package testutils

import (
	"food_ordering_backend/controllers/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPopulate(t *testing.T) {
	t.Run("should populate test users and test admins", func(t *testing.T) {
		it := assert.New(t)

		for _, u := range append(TestUsers, TestAdmins...) {
			it.NotZero(u.Email)
			it.NotZero(u.PasswordHash)
		}
	})
}

func TestLogin(t *testing.T) {
	t.Run("should send json of provided dto to login route", func(t *testing.T) {
		SetupUsersDB(t)
		it := assert.New(t)
		dto := user.AuthDTO{Email: "Aurore31@hotmail.com", Password: "9BNQgtcgRYSEAUv"}

		resp := Login(dto)
		cookie := FindCookieByName(resp.Result(), user.SessionCookieName)

		it.NotZero(cookie)
	})
}

func TestLoginAsRandomDTO(t *testing.T) {
	dtos := []user.AuthDTO{
		{Email: "Anya_Ernser@yahoo.com", Password: "hWPr911kMNyZWsc"},
		{Email: "Aurore31@hotmail.com", Password: "9BNQgtcgRYSEAUv"},
		{Email: "Julius.Keeling@hotmail.com", Password: "MblfRKEDRQvJvIK"},
	}
	t.Run("should pick random dto from provided slice and send a login request", func(t *testing.T) {
		SetupUsersDB(t)
		it := assert.New(t)
		dto, c := loginAsRandomDTO(t, dtos)

		it.Contains(dtos, dto)
		it.NotZero(c)
	})
}

func TestLoginAsRandomAdmin(t *testing.T) {
	t.Run("should pick random user from TestAdmins and perform a login request", func(t *testing.T) {
		SetupUsersDB(t)
		it := assert.New(t)
		dto, c := LoginAsRandomAdmin(t)

		it.Contains(TestAdminsDTOs, dto)
		it.NotZero(c)
	})
}

func TestLoginAsRandomUser(t *testing.T) {
	t.Run("should pick random user from TestUsers and perform a login request", func(t *testing.T) {
		SetupUsersDB(t)
		it := assert.New(t)
		dto, c := LoginAsRandomUser(t)

		it.Contains(TestUsersDTOs, dto)
		it.NotZero(c)
	})
}

func TestFindCookieByName(t *testing.T) {
	t.Run("should find a cookie in request with provided name", func(t *testing.T) {
		it := assert.New(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie := &http.Cookie{Name: "testcookie"}
			http.SetCookie(w, cookie)
		}))
		defer server.Close()
		resp, err := http.Get(server.URL)

		if it.NoError(err) {
			c := FindCookieByName(resp, "testcookie")
			it.NotNil(c)
		}
	})
}
