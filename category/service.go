package category

type Service struct {
	Repository *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Create(c Category) Category {
	return s.Repository.Create(c)
}

func (s *Service) Save(c Category) Category {
	return s.Repository.Save(c)
}

func (s *Service) FindAll() []Category {
	return s.Repository.FindAll()
}
