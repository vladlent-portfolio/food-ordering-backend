package user

type Service struct {
	jwtService *JWTService
	repo       *Repository
}

func ProvideService(r *Repository, jwtService *JWTService) *Service {
	return &Service{jwtService, r}
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

	session.UserID = u.ID
	session.User = u
	session.Token = s.jwtService.Generate(u.ID)

	if err := s.repo.CreateSession(session); err != nil {
		return session, err
	}

	return session, nil
}
