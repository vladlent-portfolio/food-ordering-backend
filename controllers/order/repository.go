package order

import (
	"food_ordering_backend/common"
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

// FindAll returns all orders for user with specified user ID.
// If user ID equals 0, return orders for all users instead.
// If Paginator is not nil, it returns paginated result.
func (r *Repository) FindAll(uid uint, p common.Paginator) ([]Order, error) {
	var orders []Order
	tx := r.preload()

	if uid != 0 {
		tx.Where("user_id = ?", uid)
	}

	if p != nil {
		tx.Scopes(common.WithPagination(p))
	}

	err := tx.Order("id ASC").Find(&orders).Error
	return orders, err
}

func (r *Repository) FindByID(id uint) (Order, error) {
	var order Order
	err := r.preload().First(&order, id).Error
	return order, err
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

// CountAll returns total amount of orders for user with specified ID.
// If ID equals 0, returns total amount of orders instead.
func (r *Repository) CountAll(uid uint) int {
	var count int64
	tx := r.db.Model(&Order{})

	if uid != 0 {
		tx.Where("user_id = ?", uid)
	}

	tx.Count(&count)
	return int(count)
}

func (r *Repository) preload() *gorm.DB {
	return r.db.
		Preload("Items").
		Preload("Items.Dish").
		Preload("Items.Dish.Category").
		Joins("User")
}
