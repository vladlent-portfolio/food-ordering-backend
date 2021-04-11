package category

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(c Category) Category {
	r.DB.Create(&c)
	return c
}

func (r *Repository) Save(c Category) Category {
	r.DB.Save(&c)
	return c
}

func (r *Repository) FindAll() []Category {
	var categories []Category
	r.DB.Find(&categories)
	return categories
}
