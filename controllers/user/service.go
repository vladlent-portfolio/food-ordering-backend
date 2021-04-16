package user

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

//func (s *Service) Login(dto AuthDTO) (string, error) {
//
//}
