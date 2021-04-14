package dish

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(d Dish) (Dish, error) {
	err := r.DB.Create(&d).Error

	if err != nil {
		return d, err
	}

	err = r.preload().First(&d).Error

	return d, err
}

func (r *Repository) Save(d Dish) (Dish, error) {
	err := r.DB.Save(&d).Error

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

func (r *Repository) FindAll() []Dish {
	var dishes []Dish
	r.preload().Find(&dishes)
	return dishes
}

func (r *Repository) Delete(d Dish) (Dish, error) {
	err := r.preload().Unscoped().Delete(&d).Error
	return d, err
}

func (r *Repository) preload() *gorm.DB {
	return r.DB.Preload("Category")
}
