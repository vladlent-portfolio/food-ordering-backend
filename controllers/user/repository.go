package user

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(u User) (User, error) {
	err := r.DB.Create(&u).Error
	return u, err
}

func (r *Repository) Find(u *User) error {
	return r.DB.Find(&u).Error
}

func (r *Repository) FindAll() []User {
	var users []User
	r.DB.Find(&users)
	return users
}

func (r *Repository) CreateSession(s Session) error {
	return r.DB.Create(&s).Error
}
