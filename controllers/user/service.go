package user

import "golang.org/x/crypto/bcrypt"

type Service struct {
	Repository *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Create(dto AuthDTO) (User, error) {
	var user User
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.Email = dto.Email
	user.PasswordHash = hashedPass

	return s.Repository.Create(user)
}
