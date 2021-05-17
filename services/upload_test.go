package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpload_ParseAndSave(t *testing.T) {
	var upload Upload

	beforeEach := func() {
		upload = Upload{
			AllowedTypes: []string{"image/png", "image/webp", "image/jpeg"},
			MaxFileSize:  0,
			Root:         "",
			FormDataKey:  "file",
		}
	}

	t.Run("should allow any file type if AllowedTypes is empty", func(t *testing.T) {
		beforeEach()
		it := assert.New(t)

	})

	t.Run("should save file into provided root folder", func(t *testing.T) {
		it := assert.New(t)

	})

	t.Run("should return 400 if there is no file with provided key", func(t *testing.T) {
		it := assert.New(t)

	})

	t.Run("should return 413 if filesize is too big", func(t *testing.T) {
		it := assert.New(t)

	})

}
