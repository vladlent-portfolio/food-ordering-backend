package user

import (
	"food_ordering_backend/common"
	"gorm.io/gorm"
)

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

func (r *Repository) FindAll(p common.Paginator) []User {
	var users []User
	tx := r.db

	if p != nil {
		r.db.Scopes(common.WithPagination(p))
	}

	tx.Find(&users)
	return users
}

func (r *Repository) CreateSession(s Session) error {
	return r.db.Create(&s).Error
}

func (r *Repository) FindSessionByToken(token string) (Session, error) {
	var session Session
	err := r.db.Where("token = ?", token).Joins("User").Find(&session).Error
	return session, err
}

func (r *Repository) DeleteSession(token string) error {
	return r.db.Delete(&Session{}, "token = ?", token).Error
}

func (r *Repository) DeleteAllSessions(u User) error {
	var s Session
	return r.db.Where("user_id = ?", u.ID).Delete(&s).Error
}

func (r *Repository) CountAll() int {
	var count int64
	r.db.Model(&User{}).Count(&count)
	return int(count)
}
