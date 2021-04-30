package common

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func IsDuplicateKeyErr(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}

func RandomInt(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}

func MIMEType(file io.ReadSeeker) (string, error) {
	fileHeader := make([]byte, 512)

	if _, err := file.Read(fileHeader); err != nil {
		return "", err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	return http.DetectContentType(fileHeader), nil
}
