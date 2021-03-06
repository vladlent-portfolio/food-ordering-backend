package order

import (
	"errors"
	"fmt"
	"food_ordering_backend/common"
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

// FindAll returns all orders for user with specified user ID.
// If user ID equals 0, return orders for all users instead.
// If Paginator is not nil, it returns paginated result.
func (s *Service) FindAll(uid uint, p common.Paginator) ([]Order, error) {
	return s.repo.FindAll(uid, p)
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

func (s *Service) Create(itemsDTO []ItemCreateDTO, u user.User) (Order, error) {
	items, err := s.ItemsFromDTOs(itemsDTO)

	if err != nil {
		return Order{}, err
	}

	o := Order{
		UserID: u.ID,
		Status: StatusCreated,
		Items:  items,
		Total:  CalcTotal(items),
	}

	return s.repo.Create(o)
}

func (s *Service) Update(o Order, dto UpdateDTO) (Order, error) {
	items, err := s.ItemsFromDTOs(dto.Items)

	if err != nil {
		return Order{}, err
	}

	if err := s.repo.DeleteItemsByID(Items(o.Items).IDs()); err != nil {
		return Order{}, err
	}

	o.Status = dto.Status
	o.UserID = dto.UserID
	o.Total = dto.Total
	o.Items = items

	return s.repo.Save(o)
}

func (s *Service) UpdateStatus(id uint, status Status) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *Service) ItemsFromDTOs(itemsDTO []ItemCreateDTO) ([]Item, error) {
	ids := make([]uint, len(itemsDTO))

	for i, dto := range itemsDTO {
		ids[i] = dto.ID
	}

	dishes, err := s.dishes.FindByIDs(ids)

	if err != nil {
		return nil, err
	}

	items := make([]Item, len(itemsDTO))
	for i, dto := range itemsDTO {
		d, ok := dish.Dishes(dishes).Find(func(d dish.Dish, index int) bool {
			return d.ID == dto.ID
		})

		if !ok {
			return nil, &ErrDishID{ID: dto.ID}
		}

		item := Item{
			DishID:   dto.ID,
			Quantity: dto.Quantity,
			Dish:     d,
		}
		items[i] = item
	}

	return items, nil
}

// CountAll returns total amount of orders for user with specified ID.
// If ID equals 0, returns total amount of orders instead.
func (s *Service) CountAll(uid uint) int {
	return s.repo.CountAll(uid)
}
