package common

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
