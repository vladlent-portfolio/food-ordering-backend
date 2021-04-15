package user_test

import (
	"encoding/json"
	"food_ordering_backend/controllers/user"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var db = database.MustGetTest()
var r = router.Setup(db)

func TestAPI(t *testing.T) {
	//t.Run("GET /users", func(t *testing.T) {
	//	send := sendReq(http.MethodGet, "/users")
	//	t.Run("should return a list of users", func(t *testing.T) {
	//		setupDB(t)
	//		it := assert.New(t)
	//		resp := send("")
	//	})
	//})

	t.Run("POST /users", func(t *testing.T) {
		send := sendReq(http.MethodPost, "/users")

		t.Run("should create a user", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			resp := send(`{"email":"example@user.com", "password": "secretpass"}`)

			require.Equal(t, http.StatusCreated, resp.Code)
			var newUser user.ResponseDTO
			err := json.NewDecoder(resp.Body).Decode(&newUser)
			require.NoError(t, err)

			it.Equal("example@user.com", newUser.Email)
			it.NotZero(newUser.ID)
		})

		t.Run("should return 422 if json is invalid", func(t *testing.T) {
			it := assert.New(t)
			resp := send(`{"email":"example@user.com", "password": "secretpass`)
			it.Equal(http.StatusUnprocessableEntity, resp.Code)

			resp = send(`{"email":"example@user.com", "password": "sec`)
			it.Equal(http.StatusUnprocessableEntity, resp.Code)

			resp = send(`{"email":"example", "password": "secretpass`)
			it.Equal(http.StatusUnprocessableEntity, resp.Code)
		})

		t.Run("should return 409 if user with provided email already exists", func(t *testing.T) {
			resp := send(`{"email":"example@user.com", "password": "secretpass"}`)
			require.Equal(t, http.StatusCreated, resp.Code)

			resp = send(`{"email":"example@user.com", "password": "secretpass"}`)
			require.Equal(t, http.StatusConflict, resp.Code)
		})
	})
}

func setupDB(t *testing.T) {
	t.Cleanup(cleanup)
}

func sendReq(method, target string) func(body string) *httptest.ResponseRecorder {
	return func(body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, target, strings.NewReader(body))
		r.ServeHTTP(w, req)
		return w
	}
}

func cleanup() {
	db.Exec("TRUNCATE users;")
}
