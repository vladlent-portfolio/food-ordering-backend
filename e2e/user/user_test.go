package user_test

import (
	"encoding/json"
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"food_ordering_backend/controllers/user"
	"food_ordering_backend/database"
	"food_ordering_backend/e2e/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"

	"testing"
)

var db = database.MustGet()
var testUsers = testutils.TestUsers
var testAdmins = testutils.TestAdmins

func TestAPI(t *testing.T) {
	t.Run("GET /users", func(t *testing.T) {
		send := testutils.ReqWithCookie(http.MethodGet, "/users")
		t.Run("should return a list of users", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			it := assert.New(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := send(c, "")

			it.Equal(http.StatusOK, resp.Code)

			var dtos user.DTOsWithPagination
			err := json.NewDecoder(resp.Body).Decode(&dtos)
			require.NoError(t, err)

			it.Len(dtos.Users, len(testUsers)+len(testAdmins))
		})

		t.Run("should work with pagination", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			it := assert.New(t)

			users := make([]user.User, 0, len(testutils.TestUsers)+len(testutils.TestAdmins))
			users = append(users, testutils.TestUsers...)
			users = append(users, testutils.TestAdmins...)
			testutils.SortUsersByID(users)

			tests := []struct {
				query          string
				expectedOrders []user.User
				expectedPage   int
				expectedLimit  int
			}{
				{"limit=2", users[:2], 0, 2},
				{"limit=3&page=2", users[3:], 2, 3},
				{"page=2", nil, 2, 10},
				{"", users, 0, 10},
			}
			_, c := testutils.LoginAsRandomAdmin(t)

			for _, tc := range tests {
				resp := testutils.ReqWithCookie(http.MethodGet, "/users?"+tc.query)(c, "")

				if it.Equal(http.StatusOK, resp.Code) {
					var dto user.DTOsWithPagination

					if it.NoError(json.NewDecoder(resp.Body).Decode(&dto)) {
						it.Equal(tc.expectedPage, dto.Pagination.Page)
						it.Equal(tc.expectedLimit, dto.Pagination.Limit)
						it.Equal(len(users), dto.Pagination.Total)

						for i, u := range dto.Users {
							it.Equal(users[i].ID, u.ID)
							it.Equal(users[i].Email, u.Email)
							it.Equal(users[i].IsAdmin, u.IsAdmin)
						}
					}
				}
			}
		})

		testutils.RunAuthTests(t, http.MethodGet, "/users", true)

		t.Run("GET /users/me", func(t *testing.T) {
			send := testutils.ReqWithCookie(http.MethodGet, "/users/me")
			t.Run("should return info about authorized user", func(t *testing.T) {
				testutils.SetupUsersDB(t)
				it := assert.New(t)
				dto, c := testutils.LoginAsRandomUser(t)
				resp := send(c, "")

				if it.Equal(http.StatusOK, resp.Code) {
					var respDTO user.ResponseDTO
					err := json.NewDecoder(resp.Body).Decode(&respDTO)

					if it.NoError(err) {
						it.NotZero(respDTO.ID)
						it.Equal(respDTO.Email, dto.Email)
						it.False(respDTO.IsAdmin)
						it.False(respDTO.CreatedAt.IsZero())
					}
				}
			})

			testutils.RunAuthTests(t, http.MethodGet, "/users/me", false)
		})

		t.Run("GET /users/logout", func(t *testing.T) {
			logout := testutils.ReqWithCookie(http.MethodGet, "/users/logout")

			t.Run("should remove auth cookie and remove only current session for the user from db", func(t *testing.T) {
				testutils.SetupUsersDB(t)
				it := assert.New(t)
				dto, c := testutils.LoginAsRandomUser(t)
				totalSessions := 5

				for i := 0; i < totalSessions-1; i++ {
					resp := testutils.Login(dto)

					if it.Equal(http.StatusOK, resp.Code) {
						c = testutils.FindCookieByName(resp.Result(), user.SessionCookieName)
					}
				}

				resp := logout(c, "")

				if it.Equal(http.StatusOK, resp.Code) {
					cookie := testutils.FindCookieByName(resp.Result(), user.SessionCookieName)

					if it.NotNil(cookie) {
						it.Less(cookie.MaxAge, 0)
					}

					var sessions []user.Session

					if it.NoError(db.Joins("User").Where("email = ?", dto.Email).Find(&sessions).Error) {
						it.Len(sessions, totalSessions-1, "expected to remove only current session")

						for _, session := range sessions {
							it.NotEqual(c.Value, session.Token)
						}
					}
				}
			})

			testutils.RunAuthTests(t, http.MethodGet, "/users/me", false)
		})
	})

	t.Run("POST /users", func(t *testing.T) {
		t.Run("POST /users/signup", func(t *testing.T) {
			send := testutils.SendReq(http.MethodPost, "/users/signup")

			t.Run("should create a user and return it", func(t *testing.T) {
				testutils.SetupUsersDB(t)
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

			t.Run("should immediately authorize user", func(t *testing.T) {
				testutils.SetupUsersDB(t)
				it := assert.New(t)
				dto := user.AuthDTO{
					Email:    "example@user.com",
					Password: "secretpass",
				}

				data, err := json.Marshal(&dto)

				it.NoError(err)

				resp := send(string(data))

				if it.Equal(http.StatusCreated, resp.Code) {
					var responseDTO user.ResponseDTO
					if it.NoError(json.NewDecoder(resp.Body).Decode(&responseDTO)) {
						c := testutils.FindCookieByName(resp.Result(), user.SessionCookieName)
						validateSessionCookie(t, c)

						var session user.Session

						query := db.Where("user_id = ?", responseDTO.ID).Joins("User").First(&session)

						if it.NoError(query.Error) {
							it.Equal(dto.Email, session.User.Email)
							it.Equal(c.Value, session.Token)
						}
					}

				}
			})

			t.Run("should return 422 if json is invalid", func(t *testing.T) {
				it := assert.New(t)
				resp := send(`{"email":"example@user.com", "password": "secretpass`)
				it.Equal(http.StatusUnprocessableEntity, resp.Code)

				resp = send(`{"email":"example@user.com", "password": "sec"}`)
				it.Equal(http.StatusUnprocessableEntity, resp.Code)

				resp = send(`{"email":"example", "password": "secretpass"}`)
				it.Equal(http.StatusUnprocessableEntity, resp.Code)
			})

			t.Run("should return 409 if user with provided email already exists", func(t *testing.T) {
				testutils.SetupUsersDB(t)
				it := assert.New(t)
				resp := send(`{"email":"example@user.com", "password": "secretpass"}`)
				if it.Equal(http.StatusCreated, resp.Code) {
					resp = send(`{"email":"example@user.com", "password": "secretpass"}`)
					require.Equal(t, http.StatusConflict, resp.Code)
				}
			})
		})

		t.Run("POST /users/signin", func(t *testing.T) {
			t.Run("should sign in a user and add session to db", func(t *testing.T) {
				testutils.SetupUsersDB(t)
				it := assert.New(t)

				for i, dto := range testutils.TestUsersDTOs {
					resp := testutils.Login(dto)

					if it.Equal(http.StatusOK, resp.Code) {
						var responseDTO user.ResponseDTO
						if it.NoError(json.NewDecoder(resp.Body).Decode(&responseDTO)) {
							it.Equal(dto.Email, responseDTO.Email)
							it.NotZero(responseDTO.ID)
							it.NotZero(responseDTO.CreatedAt)
							it.False(responseDTO.IsAdmin)
						}

						c := testutils.FindCookieByName(resp.Result(), user.SessionCookieName)
						validateSessionCookie(t, c)

						var sessions []user.Session
						db.Find(&sessions)

						if it.Len(sessions, i+1) {
							it.Equal(sessions[i].Token, c.Value)
							it.False(sessions[i].CreatedAt.IsZero())
						}
					}

				}

			})

			t.Run("should return 404 if provided email or password is incorrect", func(t *testing.T) {
				testutils.SetupUsersDB(t)
				it := assert.New(t)
				u1 := testutils.TestUsersDTOs[common.RandomInt(len(testutils.TestUsersDTOs))]
				u2 := testutils.TestUsersDTOs[common.RandomInt(len(testutils.TestUsersDTOs))]

				u1.Email = "email@not.exist"

				resp := testutils.Login(u1)
				it.Equal(http.StatusNotFound, resp.Code)

				u2.Password = "some-random-pass"
				resp = testutils.Login(u2)
				it.Equal(http.StatusNotFound, resp.Code)
			})
		})
	})

}

func validateSessionCookie(t *testing.T, c *http.Cookie) {
	it := assert.New(t)
	if it.NotNil(c) {
		it.NotZero(c.Value)
		it.True(c.HttpOnly)
		it.Equal(http.SameSiteLaxMode, c.SameSite)
		it.Equal("/", c.Path)
		it.Equal(0, c.MaxAge)
		it.True(c.Secure, "expected cookie to be Secure")
		it.Equal(config.ClientURL.Hostname(), c.Domain, "expected cookie to have a proper Domain")
	}
}
