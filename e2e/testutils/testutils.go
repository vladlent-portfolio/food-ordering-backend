package testutils

import (
	"bytes"
	"food_ordering_backend/config"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
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

func UploadReqWithCookie(method, target, formField string) func(c *http.Cookie, fileName string, file io.Reader) *httptest.ResponseRecorder {
	return func(c *http.Cookie, fileName string, file io.Reader) *httptest.ResponseRecorder {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		recorder := httptest.NewRecorder()

		fw, err := writer.CreateFormFile(formField, fileName)
		noError(err)
		_, err = io.Copy(fw, file)
		noError(err)
		writer.Close()

		req := httptest.NewRequest(method, target, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		if c != nil {
			req.AddCookie(c)
		}

		Router.ServeHTTP(recorder, req)
		return recorder
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

// PathToFile combines absolute file path to testutils folder with provided path.
// Useful for testing, since all test executables are created in OS's temp folder.
func PathToFile(path string) string {
	// Test executables are created in temp dir, so we need a reference to a current file.
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, path)
}

func CleanupStaticFolder() {
	if err := os.RemoveAll(config.StaticDirAbs); err != nil {
		panic(err)
	}
}

func noError(err error) {
	if err != nil {
		panic(err)
	}
}
