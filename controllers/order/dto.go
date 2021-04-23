package order

import (
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/user"
	"time"
)

type DTO struct {
	ID        uint      `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StatusID  uint      `json:"status_id,omitempty"`
	Status    Status    `json:"status"`
	UserID    uint      `json:"user_id,omitempty"`
	User      user.User `json:"user"`
	Total     float64   `json:"total" binding:"gt=0"`
	Items     []Item    `json:"items" binding:"required,gt=0"`
}

type ItemDTO struct {
	ID       uint `json:"id,omitempty"`
	OrderID  uint `json:"order_id" binding:"required"`
	DishID   uint `json:"dish_id" binding:"required"`
	Dish     dish.Dish
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type StatusDTO struct {
	ID    uint   `json:"id" binding:"required"`
	Title string `json:"title" binding:"required,min=2,max=30"`
}
