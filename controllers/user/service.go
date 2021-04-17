package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
)

const JWTSecret = "secret string"
const JWTUserIDKey = "uid"

type AuthClaims struct {
	UserID uint `json:"uid"`
}

func (c AuthClaims) Valid(v *jwt.ValidationHelper) error {
	if c.UserID == 0 {
		return errors.New("invalid user id")
	}

	return nil
}

type Service struct {
	repo *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Create(dto AuthDTO) (User, error) {
	user, err := CreateFromDTO(dto)

	if err != nil {
		return user, err
	}

	return s.repo.Create(user)
}

func (s *Service) FindAll() []User {
	return s.repo.FindAll()
}

func (s *Service) FindByID(id uint) (User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Login(dto AuthDTO) (Session, error) {
	var session Session
	u, err := s.repo.FindByEmail(dto.Email)

	if err != nil {
		return session, err
	}

	if err := u.ValidatePassword(dto.Password); err != nil {
		return session, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{UserID: u.ID})
	ss, _ := token.SignedString([]byte(JWTSecret))

	session.UserID = u.ID
	session.User = u
	session.Token = ss

	if err := s.repo.CreateSession(session); err != nil {
		return session, err
	}

	return session, nil
}
