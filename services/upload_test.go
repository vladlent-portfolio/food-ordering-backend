package services_test

import (
	"food_ordering_backend/e2e/testutils"
	"food_ordering_backend/services"
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
var upload *services.Upload
var router *gin.Engine
var defaultHandle = func(c *gin.Context) { upload.ParseAndSave(c, fileName) }

func beforeEach(handle gin.HandlerFunc) {
	upload = &services.Upload{
		AllowedTypes: []string{"text/plain; charset=utf-8", "image/png", "image/webp", "image/jpeg"},
		MaxFileSize:  10000,
		Root:         os.TempDir(),
		FormDataKey:  "file",
	}

	router = setupRouter(handle)
}

func TestUpload_ParseAndSave(t *testing.T) {
	t.Run("should allow any file type if AllowedTypes is empty", func(t *testing.T) {
		beforeEach(defaultHandle)
		upload.AllowedTypes = []string{}
		it := assert.New(t)
		files := []io.Reader{testutils.CreateTextFile(100), testutils.CreateImagePNG(50, 50), testutils.CreateImageJPEG(35, 35)}

		for _, file := range files {
			resp := sendFile(upload.FormDataKey, file)
			it.Equal(http.StatusOK, resp.Code)
		}
	})

	t.Run("should save file into provided root folder with correct extension", func(t *testing.T) {
		beforeEach(defaultHandle)
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
		beforeEach(defaultHandle)
		file := testutils.CreateTextFile(50)

		resp := sendFile("random-key", file)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 413 if filesize is too big", func(t *testing.T) {
		beforeEach(defaultHandle)
		upload.MaxFileSize = 1
		file := testutils.CreateTextFile(50)

		resp := sendFile(upload.FormDataKey, file)
		assert.Equal(t, http.StatusRequestEntityTooLarge, resp.Code)
	})

	t.Run("should return 415 if provided file type isn't supported", func(t *testing.T) {
		beforeEach(defaultHandle)
		upload.AllowedTypes = []string{"text/plain"}
		files := []io.Reader{testutils.CreateImagePNG(50, 50), testutils.CreateImageJPEG(35, 35)}

		for _, file := range files {
			resp := sendFile(upload.FormDataKey, file)
			assert.Equal(t, http.StatusUnsupportedMediaType, resp.Code)
		}
	})

	t.Run("should return absolute path to saved file on success", func(t *testing.T) {
		it := assert.New(t)
		handle := func(c *gin.Context) {
			p := upload.ParseAndSave(c, "document")
			it.Equal(filepath.Join(upload.Root, "document.txt"), p)
		}
		beforeEach(handle)

		sendFile(upload.FormDataKey, testutils.CreateTextFile(50))
	})

	t.Run("should return empty strings on error", func(t *testing.T) {
		it := assert.New(t)
		handle := func(c *gin.Context) {
			p := upload.ParseAndSave(c, "document")
			it.Equal("", p)
		}
		beforeEach(handle)
		upload.AllowedTypes = []string{"image/png"}

		sendFile(upload.FormDataKey, testutils.CreateTextFile(50))
	})
}

func TestUpload_AllowedType(t *testing.T) {
	u := services.Upload{AllowedTypes: []string{"image/png", "image/jpeg"}}
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

func setupRouter(handle gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.POST("/upload", handle)
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
