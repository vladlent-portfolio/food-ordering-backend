package category

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
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
