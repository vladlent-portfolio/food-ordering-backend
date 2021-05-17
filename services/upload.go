package services

import (
	"github.com/gin-gonic/gin"
)

type Upload struct {
	// AllowedTypes is a slice of supported MIME-Types. If type of uploaded
	// file isn't in this list, the request will be rejected with http.StatusUnsupportedMediaType
	// Empty list means that any MIME-Type is allowed.
	AllowedTypes []string

	// MaxFileSize sets a limit for the size of uploaded file. If the file has
	// bigger size, the request will be rejected with http.StatusRequestEntityTooLarge.
	MaxFileSize int64

	// Root should be an absolute path to a directory which will be used as a root
	// for all saved files. This property is ignored if set to an empty string.
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
func (s *Upload) ParseAndSave(c *gin.Context, name string) {

}
