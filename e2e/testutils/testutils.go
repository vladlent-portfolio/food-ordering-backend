package testutils

import (
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var db = database.MustGetTest()
var Router = router.Setup(db)

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

func EqualTimestamps(t1, t2 time.Time) bool {
	// There can be slight difference between cached user and user from db
	// so we compare string representation instead
	return t1.Format(time.RFC3339) == t2.Format(time.RFC3339)
}
