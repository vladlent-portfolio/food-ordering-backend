package services

import (
	"bytes"
	"food_ordering_backend/e2e/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var fileName string
var upload Upload
var router *gin.Engine

func beforeEach() {
	upload = Upload{
		AllowedTypes: []string{"image/png", "image/webp", "image/jpeg"},
		MaxFileSize:  0,
		Root:         "",
		FormDataKey:  "file",
	}

	fileName = "random-file"
	router = setupRouter()
}

func TestUpload_ParseAndSave(t *testing.T) {
	t.Run("should allow any file type if AllowedTypes is empty", func(t *testing.T) {
		beforeEach()
		it := assert.New(t)

	})

	t.Run("should save file into provided root folder", func(t *testing.T) {
		it := assert.New(t)

	})

	t.Run("should return 400 if there is no file with provided key", func(t *testing.T) {
		it := assert.New(t)

	})

	t.Run("should return 413 if filesize is too big", func(t *testing.T) {
		it := assert.New(t)

	})

}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.POST("/upload", func(c *gin.Context) {
		upload.ParseAndSave(c, fileName)
	})
	return r
}

func sendFile(fieldName string, file io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	form, body := testutils.MultipartWithFile(fileName, file)
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", form.FormDataContentType())
	router.ServeHTTP(w, req)
	return w
}
