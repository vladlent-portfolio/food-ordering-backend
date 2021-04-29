package category

import (
	"os"
	"path/filepath"
)

type Service struct {
	repo *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
}

func PathToImages() string {
	main, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Join(main, "/static/categories")
}

func (s *Service) Create(c Category) (Category, error) {
	return s.repo.Create(c)
}

func (s *Service) Save(c Category) (Category, error) {
	return s.repo.Save(c)
}

func (s *Service) FindByID(id uint) (Category, error) {
	return s.repo.FindByID(id)
}

func (s *Service) FindAll() []Category {
	return s.repo.FindAll()
}

func (s *Service) Delete(c Category) (Category, error) {
	return s.repo.Delete(c)
}
