package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash []byte `gorm:"not null"`
	IsAdmin      bool   `gorm:"default:false"`
}

type Session struct {
	Token     string `gorm:"primaryKey"`
	UserID    uint
	CreatedAt time.Time
	User      User `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

var ErrInvalidPassword = errors.New("invalid password")

func (u *User) SetPassword(password string) {
	// Basically, an error can only occur if we provide an invalid cost so we ignore it.
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	u.PasswordHash = hashedPass
}

func (u *User) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return ErrInvalidPassword
	}
	return nil
}
