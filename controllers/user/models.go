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

func (u *User) SetPassword(password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		return err
	}

	u.PasswordHash = hashedPass

	return nil
}

//func (u *User) ValidatePassword(password string) bool {
//
//}
