package testutils

import (
	"encoding/json"
	"food_ordering_backend/common"
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
var Router = router.Setup(db)

var TestUsersDTOs = []user.AuthDTO{
	{Email: "Kallie_Larson@hotmail.com", Password: "_GIGcnAkjjsbkzk"},
	{Email: "Hellen_Bogan26@hotmail.com", Password: "sgDOB7qIseBkpS3"},
	{Email: "Stella.Wolff@yahoo.com", Password: "kn_yt5XoDIexljw"},
}
var TestAdminsDTOs = []user.AuthDTO{
	{Email: "Anya_Ernser@yahoo.com", Password: "hWPr911kMNyZWsc"},
	{Email: "Aurore31@hotmail.com", Password: "9BNQgtcgRYSEAUv"},
	{Email: "Julius.Keeling@hotmail.com", Password: "MblfRKEDRQvJvIK"},
}
var TestUsers = make([]user.User, 3)
var TestAdmins = make([]user.User, 3)

func init() {
	populateTestUsers()
	populateTestAdmins()
}

func populateTestUsers() {
	for i, dto := range TestUsersDTOs {
		u := user.CreateFromDTO(dto)

		TestUsers[i] = u
	}
}

func populateTestAdmins() {
	for i, dto := range TestAdminsDTOs {
		u := user.CreateFromDTO(dto)

		u.IsAdmin = true
		TestAdmins[i] = u
	}
}

// SetupUsersDB inserts TestUsers and TestAdmins into db.
// Users and sessions tables will be truncated before and after test run.
func SetupUsersDB(t *testing.T) {
	cleanup := func() {
		db.Exec("TRUNCATE users, sessions CASCADE;")
	}

	t.Cleanup(cleanup)
	cleanup()

	require.NoError(t, db.Create(&TestUsers).Error)
	require.NoError(t, db.Create(&TestAdmins).Error)
}

func RunAuthTests(t *testing.T, method, target string, adminOnly bool) {
	t.Run("should return 401 user is unauthorized", func(t *testing.T) {
		t.Run("should return 401 if there is no session cookie in the request", func(t *testing.T) {
			resp := SendReq(method, target)("")
			assert.Equal(t, http.StatusUnauthorized, resp.Code)
		})
	})

	if adminOnly {
		t.Run("should return 401 if user is not admin", func(t *testing.T) {
			SetupUsersDB(t)
			_, c := LoginAsRandomUser(t)
			resp := ReqWithCookie(method, target)(c, "")
			assert.Equal(t, http.StatusUnauthorized, resp.Code)
		})
	}
}

func LoginAsRandomUser(t *testing.T) (user.AuthDTO, *http.Cookie) {
	return loginAsRandomDTO(t, TestUsersDTOs)
}

func LoginAsRandomAdmin(t *testing.T) (user.AuthDTO, *http.Cookie) {
	return loginAsRandomDTO(t, TestAdminsDTOs)
}

func loginAsRandomDTO(t *testing.T, dtos []user.AuthDTO) (user.AuthDTO, *http.Cookie) {
	randomIndex := common.RandomInt(len(dtos))
	dto := dtos[randomIndex]
	resp := Login(dto)
	authCookie := FindCookieByName(resp.Result(), user.SessionCookieName)

	require.NotEmpty(t, authCookie)
	return dto, authCookie
}

func Login(dto user.AuthDTO) *httptest.ResponseRecorder {
	data, _ := json.Marshal(&dto)
	return SendReq(http.MethodPost, "/users/signin")(string(data))
}

func SendReq(method, target string) func(body string) *httptest.ResponseRecorder {
	return func(body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, target, strings.NewReader(body))
		Router.ServeHTTP(w, req)
		return w
	}
}

func ReqWithCookie(method, target string) func(c *http.Cookie, body string) *httptest.ResponseRecorder {
	return func(c *http.Cookie, body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, target, strings.NewReader(body))

		if c != nil {
			req.AddCookie(c)
		}

		Router.ServeHTTP(w, req)
		return w
	}
}

func FindCookieByName(resp *http.Response, name string) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}
