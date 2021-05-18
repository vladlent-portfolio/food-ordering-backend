package services

import (
	"food_ordering_backend/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

type Upload struct {
	// AllowedTypes is a slice of supported MIME-Types. If type of uploaded
	// file isn't in this list, the request will be rejected with http.StatusUnsupportedMediaType
	// Empty list means that any MIME-Type is allowed.
	AllowedTypes []string

	// MaxFileSize sets a limit for the size of uploaded file. If the file has
	// bigger size, the request will be rejected with http.StatusRequestEntityTooLarge.
	MaxFileSize int64

	// Root should be an absolute path to a directory which will be used
	// as a root for all saved files.
	Root string

	// FormDataKey is the name of the field in form-data for file lookup.
	FormDataKey string
}

// ParseAndSave parses the request in order to find a file for upload
// using Upload.FormDataKey. It will check file's size and MIMI-Type
// using Upload.MaxFileSize and Upload.AllowedTypes respectively.
//
// If all checks are successful, the file will be saved with provided name.
// File extension will be appended automatically from detected MIME-Type.
//
// The request will be aborted instantly, with appropriate status code,
// on the first encountered error.
//
// Returns absolute path to saved file or empty string on error.
func (s *Upload) ParseAndSave(c *gin.Context, name string) string {
	fileHeader, err := c.FormFile(s.FormDataKey)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return ""
	}

	if fileHeader.Size > s.MaxFileSize {
		c.Status(http.StatusRequestEntityTooLarge)
		return ""
	}

	file, err := fileHeader.Open()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return ""
	}

	defer file.Close()

	mimeType, err := common.MIMEType(file)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return ""
	}

	if len(s.AllowedTypes) > 0 && !s.AllowedType(mimeType) {
		c.Status(http.StatusUnsupportedMediaType)
		return ""
	}

	ext := common.ExtensionByType(mimeType)
	fPath := filepath.Join(s.Root, name+ext)

	if err := os.MkdirAll(filepath.Dir(fPath), os.ModeDir); err != nil {
		c.Status(http.StatusInternalServerError)
		return ""
	}

	if err := c.SaveUploadedFile(fileHeader, fPath); err != nil {
		c.Status(http.StatusInternalServerError)
		return ""
	}

	return fPath
}

// AllowedType checks if provided MIME-Type is in the list of allowed MIME-Types.
func (s *Upload) AllowedType(mimetype string) bool {
	for _, allowedType := range s.AllowedTypes {
		if allowedType == mimetype {
			return true
		}
	}
	return false
}
