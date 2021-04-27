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

func (r *Repository) Create(o Order) (Order, error) {
	err := r.db.Create(&o).Error

	if err != nil {
		return o, err
	}

	err = r.preload().First(&o).Error
	return o, err
}

func (r *Repository) Save(o Order) (Order, error) {
	if err := r.db.Omit("User").Save(&o).Error; err != nil {
		return Order{}, err
	}

	var updated Order
	if err := r.preload().First(&updated, o.ID).Error; err != nil {
		return Order{}, err
	}

	return updated, nil
}

func (r *Repository) FindAll() ([]Order, error) {
	var orders []Order
	err := r.preload().Find(&orders).Error
	return orders, err
}

func (r *Repository) FindByID(id uint) (Order, error) {
	var order Order
	err := r.preload().First(&order, id).Error
	return order, err
}

func (r *Repository) FindByUID(uid uint) ([]Order, error) {
	var orders []Order
	err := r.preload().Where("user_id = ?", uid).Find(&orders).Error
	return orders, err
}

func (r *Repository) UpdateStatus(id uint, status Status) error {
	o := Order{
		ID: id,
	}
	return r.db.Model(&o).Update("status", status).Error
}

func (r *Repository) DeleteItemsByID(ids []uint) error {
	return r.db.Delete(Item{}, ids).Error
}

func (r *Repository) preload() *gorm.DB {
	return r.db.
		Preload("Items").
		Preload("Items.Dish").
		Preload("Items.Dish.Category").
		Joins("User")
}
