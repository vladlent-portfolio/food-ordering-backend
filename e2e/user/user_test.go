package user_test

import (
	"encoding/json"
	"food_ordering_backend/common"
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
var testUsersDTOs = []user.AuthDTO{
	{Email: "Kallie_Larson@hotmail.com", Password: "_GIGcnAkjjsbkzk"},
	{Email: "Hellen_Bogan26@hotmail.com", Password: "sgDOB7qIseBkpS3"},
	{Email: "Stella.Wolff@yahoo.com", Password: "kn_yt5XoDIexljw"},
}
var testAdminsDTOs = []user.AuthDTO{
	{Email: "Anya_Ernser@yahoo.com", Password: "hWPr911kMNyZWsc"},
	{Email: "Aurore31@hotmail.com", Password: "9BNQgtcgRYSEAUv"},
	{Email: "Julius.Keeling@hotmail.com", Password: "MblfRKEDRQvJvIK"},
}
var testUsers = make([]user.User, 3)
var testAdmins = make([]user.User, 3)

func init() {
	populateTestUsers()
	populateTestAdmins()
}

func TestAPI(t *testing.T) {
	t.Run("GET /users", func(t *testing.T) {
		send := sendWithCookie("/users")
		t.Run("should return a list of users", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			_, c := loginAsRandomAdmin(t)
			resp := send(c)

			it.Equal(http.StatusOK, resp.Code)

			var users []user.ResponseDTO
			err := json.NewDecoder(resp.Body).Decode(&users)
			require.NoError(t, err)

			it.Len(users, len(testUsers)+len(testAdmins))
		})

		t.Run("should return 401 if user is unauthorized ", func(t *testing.T) {
			resp := send(&http.Cookie{})
			assert.Equal(t, http.StatusUnauthorized, resp.Code)
		})

		t.Run("GET /users/me", func(t *testing.T) {
			send := sendWithCookie("/users/me")
			t.Run("should return info about authorized user", func(t *testing.T) {
				setupDB(t)
				it := assert.New(t)
				dto, c := loginAsRandomUser(t)
				resp := send(c)

				if it.Equal(http.StatusOK, resp.Code) {
					var respDTO user.ResponseDTO
					err := json.NewDecoder(resp.Body).Decode(&respDTO)

					if it.NoError(err) {
						it.NotZero(respDTO.ID)
						it.Equal(respDTO.Email, dto.Email)
						it.False(respDTO.IsAdmin)
					}
				}
			})

			t.Run("should return 401 if there is no session cookie in the request", func(t *testing.T) {
				resp := sendReq(http.MethodGet, "/users/me")("")
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
			})
		})

		t.Run("GET /users/logout", func(t *testing.T) {
			logout := sendWithCookie("/users/logout")

			t.Run("should remove auth cookie and remove all session for this user from db", func(t *testing.T) {
				setupDB(t)
				it := assert.New(t)
				dto, c := loginAsRandomUser(t)
				resp := logout(c)

				if it.Equal(http.StatusOK, resp.Code) {
					cookie := findCookieByName(resp.Result(), user.SessionCookieName)

					if it.NotNil(cookie) {
						it.Less(cookie.MaxAge, 0)
					}

					var sessions []user.Session

					if it.NoError(db.Joins("User").Where("email = ?", dto.Email).Find(&sessions).Error) {
						it.Len(sessions, 0)
					}
				}
			})
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

			var u user.User
			if it.NoError(db.Last(&u).Error) {
				it.NotEmpty(u.PasswordHash)
				it.Equal(newUser.Email, u.Email)
				it.False(u.CreatedAt.IsZero())
				it.False(u.IsAdmin)
			}
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

		t.Run("POST /users/signin", func(t *testing.T) {

			t.Run("should sign in a user and add session to db", func(t *testing.T) {
				setupDB(t)
				it := assert.New(t)

				for i, dto := range testUsersDTOs {
					resp := login(dto)

					it.Equal(http.StatusOK, resp.Code)

					c := findCookieByName(resp.Result(), user.SessionCookieName)

					if assert.NotNil(t, c) {
						it.NotZero(t, c.Value)
						it.True(c.HttpOnly)
						it.Equal(http.SameSiteLaxMode, c.SameSite)
						it.Equal("/", c.Path)
						it.Equal(0, c.MaxAge)
					}

					var sessions []user.Session
					db.Find(&sessions)

					if it.Len(sessions, i+1) {
						it.Equal(sessions[i].Token, c.Value)
						it.False(sessions[i].CreatedAt.IsZero())
					}
				}

			})

			t.Run("should return 403 if provided email or password is incorrect", func(t *testing.T) {
				setupDB(t)
				it := assert.New(t)
				u1 := testUsersDTOs[common.RandomInt(len(testUsersDTOs))]
				u2 := testUsersDTOs[common.RandomInt(len(testUsersDTOs))]

				u1.Email = "email@not.exist"

				resp := login(u1)
				it.Equal(http.StatusForbidden, resp.Code)

				u2.Password = "some-random-pass"
				resp = login(u2)
				it.Equal(http.StatusForbidden, resp.Code)
			})
		})
	})

}

func setupDB(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)
	require.NoError(t, db.Create(&testUsers).Error)
	require.NoError(t, db.Create(&testAdmins).Error)
}

func login(dto user.AuthDTO) *httptest.ResponseRecorder {
	data, _ := json.Marshal(&dto)
	return sendReq(http.MethodPost, "/users/signin")(string(data))
}

func loginAsRandomUser(t *testing.T) (user.AuthDTO, *http.Cookie) {
	return loginAsRandomDTO(t, testUsersDTOs)
}

func loginAsRandomAdmin(t *testing.T) (user.AuthDTO, *http.Cookie) {
	return loginAsRandomDTO(t, testAdminsDTOs)
}

func loginAsRandomDTO(t *testing.T, dtos []user.AuthDTO) (user.AuthDTO, *http.Cookie) {
	randomIndex := common.RandomInt(len(dtos))
	dto := dtos[randomIndex]
	resp := login(dto)
	authCookie := findCookieByName(resp.Result(), user.SessionCookieName)

	require.NotEmpty(t, authCookie)
	return dto, authCookie
}

func sendReq(method, target string) func(body string) *httptest.ResponseRecorder {
	return func(body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, target, strings.NewReader(body))
		r.ServeHTTP(w, req)
		return w
	}
}

func sendWithCookie(target string) func(c *http.Cookie) *httptest.ResponseRecorder {
	return func(c *http.Cookie) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.AddCookie(c)
		r.ServeHTTP(w, req)
		return w
	}
}

func populateTestUsers() {
	for i, dto := range testUsersDTOs {
		u, err := user.CreateFromDTO(dto)

		if err != nil {
			log.Fatalln(err)
		}

		testUsers[i] = u
	}
}

func populateTestAdmins() {
	for i, dto := range testAdminsDTOs {
		u, err := user.CreateFromDTO(dto)

		if err != nil {
			log.Fatalln(err)
		}

		u.IsAdmin = true
		testAdmins[i] = u
	}
}

func findCookieByName(resp *http.Response, name string) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func cleanup() {
	db.Exec("TRUNCATE users, sessions;")
}
