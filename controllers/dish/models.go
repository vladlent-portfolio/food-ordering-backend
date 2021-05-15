package dish

import (
	"food_ordering_backend/controllers/category"
)

type Dishes []Dish

type Dish struct {
	ID         uint    `gorm:"primaryKey"`
	Title      string  `gorm:"size:100;unique;not null"`
	Price      float64 `gorm:"check:price >= 0"`
	Image      *string
	CategoryID uint
	Category   category.Category `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (dishes Dishes) Find(lookup func(d Dish, index int) bool) (Dish, bool) {
	for i, dish := range dishes {
		if lookup(dish, i) {
			return dish, true
		}
	}
	return Dish{}, false
}
