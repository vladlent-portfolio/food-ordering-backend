package category

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(c Category) (Category, error) {
	err := r.db.Create(&c).Error
	return c, err
}

func (r *Repository) Save(c Category) (Category, error) {
	err := r.db.Save(&c).Error
	return c, err
}

func (r *Repository) FindByID(id uint) (Category, error) {
	var c Category
	err := r.db.First(&c, id).Error
	return c, err
}

func (r *Repository) FindAll() []Category {
	var categories []Category
	r.db.Find(&categories)
	return categories
}

func (r *Repository) Delete(c Category) (Category, error) {
	err := r.db.Unscoped().Delete(&c).Error
	return c, err
}
