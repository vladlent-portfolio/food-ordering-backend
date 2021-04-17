package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(u User) (User, error) {
	err := r.db.Create(&u).Error
	return u, err
}

func (r *Repository) FindByID(id uint) (User, error) {
	var u User
	if err := r.db.Find(&u, id).Error; err != nil {
		return u, err
	}
	return u, nil
}

func (r *Repository) FindByEmail(email string) (User, error) {
	var u User
	err := r.db.Where("email = ?", email).Find(&u).Error
	return u, err
}

func (r *Repository) FindAll() []User {
	var users []User
	r.db.Find(&users)
	return users
}

func (r *Repository) CreateSession(s Session) error {
	return r.db.Create(&s).Error
}
