package category

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Title     string `gorm:"size:255,unique,not null"`
	Removable bool
}
