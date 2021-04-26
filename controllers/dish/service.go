package dish

type Service struct {
	repo *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Create(d Dish) (Dish, error) {
	return s.repo.Create(d)
}

func (s *Service) Save(d Dish) (Dish, error) {
	return s.repo.Save(d)
}

func (s *Service) FindByID(id uint) (Dish, error) {
	return s.repo.FindByID(id)
}

func (s *Service) FindByIDs(ids []uint) ([]Dish, error) {
	return s.repo.FindByIDs(ids)
}

func (s *Service) FindAll(cid uint) []Dish {
	return s.repo.FindAll(cid)
}

func (s *Service) Delete(d Dish) (Dish, error) {
	return s.repo.Delete(d)
}
