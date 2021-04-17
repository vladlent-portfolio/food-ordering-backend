package dish

import (
	"food_ordering_backend/controllers/category"
)

type Dish struct {
	ID         uint    `gorm:"primaryKey"`
	Title      string  `gorm:"size:100"`
	Price      float64 `gorm:"check:price >= 0"`
	CategoryID uint
	Category   category.Category `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
