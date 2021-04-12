package category

type Service struct {
	Repository *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Create(c Category) (Category, error) {
	return s.Repository.Create(c)
}

func (s *Service) Save(c Category) (Category, error) {
	return s.Repository.Save(c)
}

func (s *Service) FindByID(id uint) (Category, error) {
	return s.Repository.FindByID(id)
}

func (s *Service) FindAll() []Category {
	return s.Repository.FindAll()
}

func (s *Service) Delete(c Category) (Category, error) {
	return s.Repository.Delete(c)
}
