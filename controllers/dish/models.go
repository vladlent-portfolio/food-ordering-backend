package dish

import (
	"food_ordering_backend/controllers/category"
	"gorm.io/gorm"
)

type Dish struct {
	gorm.Model
	Title      string  `gorm:"size:100"`
	Price      float64 `gorm:"check:price >= 0"`
	CategoryID uint
	Category   category.Category `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
