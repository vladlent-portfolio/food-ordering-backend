package dish

type Service struct {
	Repository *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Create(d Dish) (Dish, error) {
	return s.Repository.Create(d)
}

func (s *Service) Save(d Dish) (Dish, error) {
	return s.Repository.Save(d)
}

func (s *Service) FindByID(id uint) (Dish, error) {
	return s.Repository.FindByID(id)
}

func (s *Service) FindAll() []Dish {
	return s.Repository.FindAll()
}

func (s *Service) Delete(d Dish) (Dish, error) {
	return s.Repository.Delete(d)
}
