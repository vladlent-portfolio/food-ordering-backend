package user_test

import (
	"encoding/json"
	"food_ordering_backend/controllers/user"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var db = database.MustGetTest()
var r = router.Setup(db)
var testUsers = make([]user.User, 3)

func init() {
	dtos := []user.AuthDTO{
		{Email: "Kallie_Larson@hotmail.com", Password: "_GIGcnAkjjsbkzk"},
		{Email: "Hellen_Bogan26@hotmail.com", Password: "sgDOB7qIseBkpS3"},
		{Email: "Stella.Wolff@yahoo.com", Password: "kn_yt5XoDIexljw"},
	}

	for i, dto := range dtos {
		u, err := user.CreateFromDTO(dto)

		if err != nil {
			log.Fatalln(err)
		}

		testUsers[i] = u
	}
}

func TestAPI(t *testing.T) {
	t.Run("GET /users", func(t *testing.T) {
		send := sendReq(http.MethodGet, "/users")
		t.Run("should return a list of users", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			resp := send("")

			it.Equal(http.StatusOK, resp.Code)

			var users []user.ResponseDTO
			err := json.NewDecoder(resp.Body).Decode(&users)
			require.NoError(t, err)

			it.Len(users, len(testUsers))

			for i, u := range users {
				it.Equal(testUsers[i].ID, u.ID)
				it.Equal(testUsers[i].Email, u.Email)
			}
		})
	})

	t.Run("POST /users", func(t *testing.T) {
		send := sendReq(http.MethodPost, "/users")

		t.Run("should create a user and return it", func(t *testing.T) {
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
	cleanup()
	t.Cleanup(cleanup)

	db.Create(&testUsers)
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
