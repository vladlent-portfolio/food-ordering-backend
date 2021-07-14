package dish

import (
	"food_ordering_backend/config"
	"food_ordering_backend/controllers/category"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

type Dishes []Dish

type Dish struct {
	ID         uint    `gorm:"primaryKey"`
	Title      string  `gorm:"size:100;unique;not null"`
	Price      float64 `gorm:"check:price >= 0"`
	Image      *string
	Removable  bool `gorm:"default:true"`
	CategoryID uint
	Category   category.Category `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (d *Dish) AfterDelete(tx *gorm.DB) (err error) {
	if d.Image != nil {
		err = os.Remove(filepath.Join(config.DishesImgDirAbs, *d.Image))
		if err != nil {
			log.Println("[Dish] Error deleting image:", err)
		}
	}
	return
}

func (dishes Dishes) Find(lookup func(d Dish, index int) bool) (Dish, bool) {
	for i, dish := range dishes {
		if lookup(dish, i) {
			return dish, true
		}
	}
	return Dish{}, false
}
