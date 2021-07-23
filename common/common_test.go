package common_test

import (
	"bytes"
	"errors"
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"food_ordering_backend/tests/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/url"
	"os"
	"runtime"
	"testing"
)

func TestIsDuplicateKeyErr(t *testing.T) {
	assert.True(t, common.IsDuplicateKeyErr(errors.New("ERROR: duplicate key value violates unique constraint \"categories_pkey\" (SQLSTATE 23505)")))
	assert.False(t, common.IsDuplicateKeyErr(errors.New("92374283uasdfj")))
}

func TestIsForeignKeyErr(t *testing.T) {
	assert.True(t, common.IsForeignKeyErr(errors.New("update or delete on table \"dishes\" violates foreign key constraint \"fk_order_items_dish\" on table \"order_items\" (SQLSTATE 23503)")))
	assert.False(t, common.IsForeignKeyErr(errors.New("92374283uasdfj")))
}

func TestMIMEType(t *testing.T) {
	t.Run("should return MIME Type for provided files", func(t *testing.T) {
		it := assert.New(t)
		_, filename, _, _ := runtime.Caller(0)

		f, err := os.Open(filename)
		require.NoError(t, err)

		tests := []struct {
			file     io.ReadSeeker
			mimetype string
		}{
			{f, "text/plain; charset=utf-8"},
			{bytes.NewReader(testutils.CreateImagePNG(30, 45).Bytes()), "image/png"},
			{bytes.NewReader(testutils.CreateImageJPEG(30, 20).Bytes()), "image/jpeg"},
		}

		for _, tc := range tests {
			mime, err := common.MIMEType(tc.file)
			it.NoError(err)
			it.Equal(tc.mimetype, mime)
		}
	})

	t.Run("should correctly return mimetype for text files with size lees than 512 bytes", func(t *testing.T) {
		it := assert.New(t)
		sizes := []int{30, 69, 345, 510}

		for _, size := range sizes {
			file := testutils.CreateTextFile(size)
			mimetype, err := common.MIMEType(bytes.NewReader(file.Bytes()))

			if it.NoError(err) {
				it.Equal("text/plain; charset=utf-8", mimetype)
			}
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

	t.Run("should resolve to absolute reference without port in prod mode", func(t *testing.T) {
		config.IsProdMode = true

		t.Cleanup(func() {
			config.IsProdMode = false
		})

		it := assert.New(t)
		tests := []struct{ ref, expected string }{
			{"/docs/cv.pdf", "https://example.com/docs/cv.pdf"},
			{"docs/cv.pdf", "https://example.com/docs/cv.pdf"},
			{"./docs/cv.pdf", "https://example.com/docs/cv.pdf"},
		}

		for _, tc := range tests {
			uri := common.HostURLResolver(tc.ref)
			it.Equal(tc.expected, uri)
		}
	})

	t.Run("should resolve to absolute reference with port in dev mode", func(t *testing.T) {
		oldHost := config.HostRaw

		config.HostRaw = config.HostRaw + ":" + "1234"
		config.HostURL, _ = url.Parse(config.HostRaw)

		t.Cleanup(func() {
			config.HostRaw = oldHost
			config.HostURL, _ = url.Parse(oldHost)
		})

		it := assert.New(t)
		tests := []struct{ ref, expected string }{
			{"/docs/cv.pdf", "https://example.com:1234/docs/cv.pdf"},
			{"docs/cv.pdf", "https://example.com:1234/docs/cv.pdf"},
			{"./docs/cv.pdf", "https://example.com:1234/docs/cv.pdf"},
		}

		for _, tc := range tests {
			uri := common.HostURLResolver(tc.ref)
			it.Equal(tc.expected, uri)
		}
	})
}

func TestExtensionByType(t *testing.T) {
	t.Run("should return appropriate file extension", func(t *testing.T) {
		tests := []struct{ mimetype, ext string }{
			{"text/plain", ".txt"},
			{"image/png", ".png"},
			{"image/jpeg", ".jpeg"},
			{"application/json", ".json"},
			{"application/pdf", ".pdf"},
		}

		for _, tc := range tests {
			assert.Equal(t, tc.ext, common.ExtensionByType(tc.mimetype))
		}
	})
}
