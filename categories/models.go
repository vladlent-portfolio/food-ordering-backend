package categories

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Title     string `json:"title" gorm:"size:255"`
	Removable bool   `json:"-"`
}
