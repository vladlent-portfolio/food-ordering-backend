package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_SetPassword(t *testing.T) {
	t.Run("should generate a hash for a password", func(t *testing.T) {
		it := assert.New(t)
		var user User
		password := "password"

		user.SetPassword(password)
		it.NotEmpty(user.PasswordHash)
		it.NotEqual(password, user.PasswordHash)
	})

	t.Run("should generate a new hash everytime", func(t *testing.T) {
		it := assert.New(t)
		var user User
		password := "password"

		user.SetPassword(password)
		oldHash := user.PasswordHash

		user.SetPassword(password)
		it.NotEqual(oldHash, user.PasswordHash)
	})
}

func TestUser_ValidatePassword(t *testing.T) {
	t.Run("should return nil if provided password is correct", func(t *testing.T) {
		it := assert.New(t)
		var user User
		password := "password"

		user.SetPassword(password)
		it.NoError(user.ValidatePassword(password))
	})

	t.Run("should return error if provided password is incorrect", func(t *testing.T) {
		it := assert.New(t)
		var user User

		user.SetPassword("pass")
		it.ErrorIs(user.ValidatePassword("notpass"), ErrInvalidPassword)
	})
}
