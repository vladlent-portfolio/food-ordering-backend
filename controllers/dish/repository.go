package dish

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(c Dish) (Dish, error) {
	err := r.DB.Create(&c).Error
	return c, err
}

func (r *Repository) Save(d Dish) (Dish, error) {
	err := r.DB.Save(&d).Error
	return d, err
}

func (r *Repository) FindByID(id uint) (Dish, error) {
	var d Dish
	err := r.DB.First(&d, id).Error
	return d, err
}

func (r *Repository) FindAll() []Dish {
	var dishes []Dish
	r.DB.Find(&dishes)
	return dishes
}

func (r *Repository) Delete(d Dish) (Dish, error) {
	err := r.DB.Delete(&d).Error
	return d, err
}
