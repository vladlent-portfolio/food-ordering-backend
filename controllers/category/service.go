package category

import (
	"food_ordering_backend/config"
	"os"
	"path/filepath"
)

type Service struct {
	repo *Repository
}

func ProvideService(r *Repository) *Service {
	return &Service{r}
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

func (s *Service) FindAllDishImages(categoryID uint) ([]string, error) {
	return s.repo.FindAllDishImages(categoryID)
}

func (s *Service) DeleteDishImages(dishImages []string) error {
	for _, image := range dishImages {
		err := os.Remove(filepath.Join(config.DishesImgDirAbs, image))

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Delete(c Category) (Category, error) {
	return s.repo.Delete(c)
}
