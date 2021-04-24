package order

type Service struct {
	repo *Repository
}

func ProvideService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) FindAll() ([]Order, error) {
	return s.repo.FindAll()
}

func (s *Service) FindByUID(uid uint) ([]Order, error) {
	return s.repo.FindByUID(uid)
}
