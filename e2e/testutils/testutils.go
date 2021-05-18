package testutils

import (
	"bytes"
	"food_ordering_backend/config"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
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
var Router *gin.Engine

func init() {
	config.StaticDir = "test-static"
	config.CategoriesImgDir = config.StaticDir + "/categories"
	config.DishesImgDir = config.StaticDir + "/dishes"
	config.StaticDirAbs = PathToFile(config.StaticDir)
	config.CategoriesImgDirAbs = PathToFile(config.CategoriesImgDir)
	config.DishesImgDirAbs = PathToFile(config.DishesImgDir)

	Router = router.Setup(db)
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

func UploadReqWithCookie(method, target, fieldName string) func(c *http.Cookie, fileName string, file io.Reader) *httptest.ResponseRecorder {
	return func(c *http.Cookie, fileName string, file io.Reader) *httptest.ResponseRecorder {
		writer, body := MultipartWithFile(fieldName, fileName, file)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(method, target, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		if c != nil {
			req.AddCookie(c)
		}

		Router.ServeHTTP(recorder, req)
		return recorder
	}
}

func MultipartWithFile(fieldName, fileName string, file io.Reader) (*multipart.Writer, *bytes.Buffer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile(fieldName, fileName)
	noError(err)
	_, err = io.Copy(fw, file)
	noError(err)
	writer.Close()
	return writer, body
}

func CreateTextFile(size int) *bytes.Buffer {
	f := &bytes.Buffer{}
	f.Grow(size)
	for i := 0; i < size; i++ {
		f.WriteRune('a')
	}
	return f
}

func CreateImage(width, height int) *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.White)
		}
	}

	return img
}

func CreateImagePNG(width, height int) *bytes.Buffer {
	img := CreateImage(width, height)
	b := &bytes.Buffer{}
	b.Grow(len(img.Pix))

	if err := png.Encode(b, img); err != nil {
		panic(err)
	}

	return b
}

func CreateImageJPEG(width, height int) *bytes.Buffer {
	img := CreateImage(width, height)
	b := &bytes.Buffer{}
	b.Grow(len(img.Pix))

	if err := jpeg.Encode(b, img, nil); err != nil {
		panic(err)
	}

	return b
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
