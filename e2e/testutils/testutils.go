package testutils

import (
	"encoding/json"
	"food_ordering_backend/common"
	"food_ordering_backend/controllers/user"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"github.com/stretchr/testify/require"
	"log"
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
		u, err := user.CreateFromDTO(dto)

		if err != nil {
			log.Fatalln(err)
		}

		TestUsers[i] = u
	}
}

func populateTestAdmins() {
	for i, dto := range TestAdminsDTOs {
		u, err := user.CreateFromDTO(dto)

		if err != nil {
			log.Fatalln(err)
		}

		u.IsAdmin = true
		TestAdmins[i] = u
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

func ReqWithCookie(target string) func(c *http.Cookie) *httptest.ResponseRecorder {
	return func(c *http.Cookie) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, target, nil)

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
