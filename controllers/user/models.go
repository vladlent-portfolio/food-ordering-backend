package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash []byte `gorm:"not null"`
}

type Session struct {
	Token  string `gorm:"primaryKey"`
	UserID uint
	User   User `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (u *User) SetPassword(password string) {
	// Basically, an error can only occur if we provide an invalid cost so we ignore it.
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	u.PasswordHash = hashedPass
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	return err != nil
}
