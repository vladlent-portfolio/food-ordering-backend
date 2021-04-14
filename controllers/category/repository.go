package category

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(c Category) (Category, error) {
	err := r.DB.Create(&c).Error
	return c, err
}

func (r *Repository) Save(c Category) (Category, error) {
	err := r.DB.Save(&c).Error
	return c, err
}

func (r *Repository) FindByID(id uint) (Category, error) {
	var c Category
	err := r.DB.First(&c, id).Error
	return c, err
}

func (r *Repository) FindAll() []Category {
	var categories []Category
	r.DB.Find(&categories)
	return categories
}

func (r *Repository) Delete(c Category) (Category, error) {
	err := r.DB.Unscoped().Delete(&c).Error
	return c, err
}
