package testutils

import (
	"encoding/json"
	"food_ordering_backend/common"
	"food_ordering_backend/controllers/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

var startID uint = 1
var TestUsers = populateTestUsers(startID)
var TestAdmins = populateTestAdmins(startID + uint(len(TestUsers)))

func populateTestUsers(startID uint) []user.User {
	users := make([]user.User, 3)
	for i, dto := range TestUsersDTOs {
		u := user.CreateFromDTO(dto)
		u.ID = startID

		startID++

		users[i] = u
	}
	return users
}

func populateTestAdmins(startID uint) []user.User {
	admins := make([]user.User, 3)
	for i, dto := range TestAdminsDTOs {
		u := user.CreateFromDTO(dto)
		u.ID = startID

		startID++

		u.IsAdmin = true
		admins[i] = u
	}
	return admins
}

// SetupUsersDB inserts TestUsers and TestAdmins into db.
// Users and sessions tables will be truncated before and after test run.
func SetupUsersDB(t *testing.T) {
	req := require.New(t)
	cleanup := func() {
		req.NoError(db.Exec("TRUNCATE users, sessions CASCADE;").Error)
	}

	t.Cleanup(cleanup)
	cleanup()

	req.NoError(db.Create(&TestUsers).Error)
	req.NoError(db.Create(&TestAdmins).Error)
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

func Login(dto user.AuthDTO) *httptest.ResponseRecorder {
	data, _ := json.Marshal(&dto)
	return SendReq(http.MethodPost, "/users/signin")(string(data))
}

func LoginAs(t *testing.T, dto user.AuthDTO) *http.Cookie {
	resp := Login(dto)
	authCookie := FindCookieByName(resp.Result(), user.SessionCookieName)

	require.NotEmpty(t, authCookie)
	return authCookie
}

func LoginAsRandomUser(t *testing.T) (user.AuthDTO, *http.Cookie) {
	return loginAsRandomDTO(t, TestUsersDTOs)
}

func LoginAsRandomAdmin(t *testing.T) (user.AuthDTO, *http.Cookie) {
	return loginAsRandomDTO(t, TestAdminsDTOs)
}

func UsersEqual(t *testing.T, u1, u2 user.User) {
	it := assert.New(t)
	it.Equal(u1.ID, u2.ID)
	it.Equal(u1.Email, u2.Email)
	// TODO: Fix user times not being equal
	it.True(u1.CreatedAt.Equal(u2.CreatedAt))
	it.Equal(u1.IsAdmin, u2.IsAdmin)
	it.Equal(u1.PasswordHash, u2.PasswordHash)
}

func loginAsRandomDTO(t *testing.T, dtos []user.AuthDTO) (user.AuthDTO, *http.Cookie) {
	randomIndex := common.RandomInt(len(dtos))
	dto := dtos[randomIndex]
	return dto, LoginAs(t, dto)
}
