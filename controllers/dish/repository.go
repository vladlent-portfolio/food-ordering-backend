package dish

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(d Dish) (Dish, error) {
	err := r.db.Create(&d).Error

	if err != nil {
		return d, err
	}

	err = r.preload().First(&d).Error

	return d, err
}

func (r *Repository) Save(d Dish) (Dish, error) {
	err := r.preload().Save(&d).Error

	if err != nil {
		return d, err
	}

	err = r.preload().First(&d).Error
	return d, err
}

func (r *Repository) FindByID(id uint) (Dish, error) {
	var d Dish
	err := r.preload().First(&d, id).Error
	return d, err
}

func (r *Repository) FindByIDs(ids []uint) ([]Dish, error) {
	var dishes []Dish
	err := r.db.Find(&dishes, ids).Error
	return dishes, err
}

func (r *Repository) FindAll(cid uint) []Dish {
	var dishes []Dish

	if cid == 0 {
		r.preload().Find(&dishes)
	} else {
		r.preload().Where("category_id = ?", cid).Find(&dishes)
	}

	return dishes
}

func (r *Repository) Delete(d Dish) (Dish, error) {
	err := r.preload().Delete(&d).Error
	return d, err
}

func (r *Repository) preload() *gorm.DB {
	return r.db.Joins("Category")
}
