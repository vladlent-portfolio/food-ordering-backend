package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash []byte `gorm:"not null"`
}

type Session struct {
	Token  string `gorm:"primaryKey"`
	UserID uint
	User   User
}

//func (u *User) ValidatePassword(password string) bool {
//
//}
