package order

import (
	"errors"
	"fmt"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/user"
	"gorm.io/gorm"
)

type Service struct {
	repo   *Repository
	dishes *dish.Service
}

type ErrDishID struct {
	ID uint
}

type ErrOrderID struct {
	ID uint
}

func (e *ErrDishID) Error() string {
	return fmt.Sprintf("Dish with id %d doesn't exist", e.ID)
}

func (e *ErrOrderID) Error() string {
	return fmt.Sprintf("Order with id %d doesn't exist", e.ID)
}

func ProvideService(repo *Repository, dishes *dish.Service) *Service {
	return &Service{repo, dishes}
}

func (s *Service) FindAll() ([]Order, error) {
	return s.repo.FindAll()
}

func (s *Service) FindByID(id uint) (Order, error) {
	o, err := s.repo.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Order{}, &ErrOrderID{ID: id}
		}

		return Order{}, err
	}

	return o, nil
}

func (s *Service) FindByUID(uid uint) ([]Order, error) {
	return s.repo.FindByUID(uid)
}

func (s *Service) Create(itemsDTO []ItemRequestDTO, u user.User) (Order, error) {
	ids := make([]uint, len(itemsDTO))

	for i, dto := range itemsDTO {
		ids[i] = dto.ID
	}

	dishes, err := s.dishes.FindByIDs(ids)

	if err != nil {
		return Order{}, err
	}

	o := Order{
		UserID: u.ID,
		Status: StatusCreated,
		Items:  make([]Item, len(itemsDTO)),
	}

	for i, dto := range itemsDTO {
		d, ok := dish.Dishes(dishes).Find(func(d dish.Dish, index int) bool {
			return d.ID == dto.ID
		})

		if !ok {
			return Order{}, &ErrDishID{ID: dto.ID}
		}

		item := Item{
			DishID:   dto.ID,
			Quantity: dto.Quantity,
			Dish:     d,
		}
		o.Items[i] = item
	}

	o.Total = CalcTotal(o.Items)
	return s.repo.Create(o)
}

func (s *Service) UpdateStatus(id uint, status Status) error {
	return s.repo.UpdateStatus(id, status)
}
