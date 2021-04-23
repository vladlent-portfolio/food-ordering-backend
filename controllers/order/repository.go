package order

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}
