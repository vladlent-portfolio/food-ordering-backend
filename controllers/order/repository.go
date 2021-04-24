package order

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) FindAll() ([]Order, error) {
	var orders []Order
	err := r.preload().Find(&orders).Error
	return orders, err
}

func (r *Repository) FindByUID(uid uint) ([]Order, error) {
	var orders []Order
	err := r.preload().Where("user_id = ?", uid).Find(&orders).Error
	return orders, err
}

func (r *Repository) preload() *gorm.DB {
	return r.db.
		Preload("Items").
		Preload("Items.Dish").
		Preload("Items.Dish.Category").
		Joins("User")
}
