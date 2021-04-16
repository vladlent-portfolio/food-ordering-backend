package user

import (
	"github.com/dgrijalva/jwt-go/v4"
)

const JWTSecret = "secret string"

type Service struct {
	Repository *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Create(dto AuthDTO) (User, error) {
	user, err := CreateFromDTO(dto)

	if err != nil {
		return user, err
	}

	return s.Repository.Create(user)
}

func (s *Service) FindAll() []User {
	return s.Repository.FindAll()
}

func (s *Service) Login(dto AuthDTO) (Session, error) {
	var session Session
	u, err := s.Repository.FindByEmail(dto.Email)

	if err != nil {
		return session, err
	}

	if err := u.ValidatePassword(dto.Password); err != nil {
		return session, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"uid": u.ID})
	ss, _ := token.SignedString([]byte(JWTSecret))

	session.UserID = u.ID
	session.User = u
	session.Token = ss

	if err := s.Repository.CreateSession(session); err != nil {
		return session, err
	}

	return session, nil
}
