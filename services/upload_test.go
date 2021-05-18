package services

import (
	"food_ordering_backend/e2e/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var fileName = "test-file"
var upload *Upload
var router *gin.Engine

func beforeEach() {
	upload = &Upload{
		AllowedTypes: []string{"text/plain; charset=utf-8", "image/png", "image/webp", "image/jpeg"},
		MaxFileSize:  10000,
		Root:         os.TempDir(),
		FormDataKey:  "file",
	}

	router = setupRouter()
}

func TestUpload_ParseAndSave(t *testing.T) {
	t.Run("should allow any file type if AllowedTypes is empty", func(t *testing.T) {
		beforeEach()
		upload.AllowedTypes = []string{}
		it := assert.New(t)
		files := []io.Reader{testutils.CreateTextFile(100), testutils.CreateImagePNG(50, 50), testutils.CreateImageJPEG(35, 35)}

		for _, file := range files {
			resp := sendFile(upload.FormDataKey, file)
			it.Equal(http.StatusOK, resp.Code)
		}
	})

	t.Run("should save file into provided root folder with correct extension", func(t *testing.T) {
		beforeEach()
		it := assert.New(t)
		tests := []struct {
			file io.Reader
			ext  string
		}{
			{testutils.CreateTextFile(50), "txt"},
			{testutils.CreateImagePNG(30, 30), "png"},
			{testutils.CreateImageJPEG(30, 30), "jpeg"},
		}

		for _, tc := range tests {
			resp := sendFile(upload.FormDataKey, tc.file)

			if it.Equal(http.StatusOK, resp.Code) {
				it.FileExists(filepath.Join(upload.Root, fileName+"."+tc.ext))
			}
		}
	})

	t.Run("should return 400 if there is no file with provided key", func(t *testing.T) {
		beforeEach()
		file := testutils.CreateTextFile(50)

		resp := sendFile("random-key", file)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 413 if filesize is too big", func(t *testing.T) {
		beforeEach()
		upload.MaxFileSize = 1
		file := testutils.CreateTextFile(50)

		resp := sendFile(upload.FormDataKey, file)
		assert.Equal(t, http.StatusRequestEntityTooLarge, resp.Code)
	})

	t.Run("should return 415 if provided file type isn't supported", func(t *testing.T) {
		beforeEach()
		upload.AllowedTypes = []string{"text/plain"}
		files := []io.Reader{testutils.CreateImagePNG(50, 50), testutils.CreateImageJPEG(35, 35)}

		for _, file := range files {
			resp := sendFile(upload.FormDataKey, file)
			assert.Equal(t, http.StatusUnsupportedMediaType, resp.Code)
		}
	})
}

func TestUpload_AllowedType(t *testing.T) {
	u := Upload{AllowedTypes: []string{"image/png", "image/jpeg"}}
	t.Run("should return false if provided mimetype isn't in the list", func(t *testing.T) {
		types := []string{"text/plain", "application/json"}

		for _, mime := range types {
			assert.False(t, u.AllowedType(mime))
		}
	})

	t.Run("should return true if provided mimetype is in the list", func(t *testing.T) {
		types := []string{"image/png", "image/jpeg"}

		for _, mime := range types {
			assert.True(t, u.AllowedType(mime))
		}
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
	form, body := testutils.MultipartWithFile(fieldName, "filename", file)
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", form.FormDataContentType())
	router.ServeHTTP(w, req)
	return w
}
