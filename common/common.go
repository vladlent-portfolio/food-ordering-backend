package common

import (
	"food_ordering_backend/config"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func IsDuplicateKeyErr(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}

func RandomInt(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}

// MIMEType returns MIME type for provided file using http.DetectContentType.
func MIMEType(file io.ReadSeeker) (string, error) {
	fileHeader := make([]byte, 0, 512)

	if _, err := file.Read(fileHeader); err != nil {
		return "", err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	return http.DetectContentType(fileHeader), nil
}

// HostURLResolver resolves reference from config.HostURL to provided relative path.
// Returns absolute URL.
func HostURLResolver(relativePath string) string {
	u, _ := url.Parse(relativePath)
	return config.HostURL.ResolveReference(u).String()
}
