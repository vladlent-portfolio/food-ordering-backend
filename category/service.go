package category

type Service struct {
	Repository *Repository
}

func (s *Service) Save(c Category) Category {
	return s.Repository.Save(c)
}

func (s *Service) FindAll() []Category {
	return s.Repository.FindAll()
}
