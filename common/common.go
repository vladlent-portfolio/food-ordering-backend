package common

import (
	"food_ordering_backend/config"
	"io"
	"math/rand"
	"mime"
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
	fileHeader := make([]byte, 512)

	n, err := file.Read(fileHeader)

	if err != nil {
		return "", err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	// Explanation of why we need to slice the file header.
	// https://gist.github.com/rayrutjes/db9b9ea8e02255d62ce2#gistcomment-3418419
	return http.DetectContentType(fileHeader[:n]), nil
}

func ExtensionByType(mimeType string) string {
	contains := func(s string) bool {
		return strings.Contains(mimeType, s)
	}

	switch {
	case contains("text/plain"):
		return ".txt"
	case contains("image/jpeg"):
		return ".jpeg"
	default:
		exts, _ := mime.ExtensionsByType(mimeType)
		return exts[0]
	}
}

// HostURLResolver resolves reference from config.HostURL to provided relative path.
// Returns absolute URL.
func HostURLResolver(relativePath string) string {
	u, _ := url.Parse(relativePath)
	return config.HostURL.ResolveReference(u).String()
}
