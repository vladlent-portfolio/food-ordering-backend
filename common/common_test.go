package common

import (
	"errors"
	"food_ordering_backend/config"
	"github.com/stretchr/testify/assert"
	"net/url"
	"os"
	"runtime"
	"testing"
)

func TestIsDuplicateKeyErr(t *testing.T) {
	assert.True(t, IsDuplicateKeyErr(errors.New("ERROR: duplicate key value violates unique constraint \"categories_pkey\" (SQLSTATE 23505)")))
	assert.False(t, IsDuplicateKeyErr(errors.New("92374283uasdfj")))
}

func TestMIMEType(t *testing.T) {
	t.Run("should return MIME Type for provided file", func(t *testing.T) {
		it := assert.New(t)
		_, filename, _, _ := runtime.Caller(0)

		f, err := os.Open(filename)

		if it.NoError(err) {
			mime, err := MIMEType(f)
			it.NoError(err)
			it.Contains(mime, "text/plain")
		}
	})
}

func TestHostURLResolver(t *testing.T) {
	oldHost := config.HostRaw
	config.HostRaw = "https://example.com"
	config.HostURL, _ = url.Parse(config.HostRaw)

	t.Cleanup(func() {
		config.HostRaw = oldHost
		config.HostURL, _ = url.Parse(oldHost)
	})

	t.Run("should resolve to absolute reference", func(t *testing.T) {
		it := assert.New(t)
		tests := []struct{ ref, expected string }{
			{"/docs/cv.pdf", "https://example.com/docs/cv.pdf"},
			{"docs/cv.pdf", "https://example.com/docs/cv.pdf"},
			{"./docs/cv.pdf", "https://example.com/docs/cv.pdf"},
		}

		for _, tc := range tests {
			uri := HostURLResolver(tc.ref)
			it.Equal(tc.expected, uri)
		}
	})
}
