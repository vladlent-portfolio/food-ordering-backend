package order

import (
	"food_ordering_backend/common"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/user"
	"time"
)

type CreateDTO struct {
	Items []ItemCreateDTO `json:"items" binding:"required,gt=0,dive"`
}

type ItemCreateDTO struct {
	ID       uint `json:"id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,gt=0"`
}

type ResponseDTO struct {
	ID        uint              `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Status    Status            `json:"status"`
	UserID    uint              `json:"user_id"`
	User      user.ResponseDTO  `json:"user"`
	Total     float64           `json:"total"`
	Items     []ItemResponseDTO `json:"items"`
}

type DTOsWithPagination struct {
	Orders     []ResponseDTO        `json:"orders"`
	Pagination common.PaginationDTO `json:"pagination"`
}

type UpdateDTO struct {
	Status Status          `json:"status" binding:"required,min=0,max=3"`
	UserID uint            `json:"user_id" binding:"required"`
	Total  float64         `json:"total" binding:"required,min=0"`
	Items  []ItemCreateDTO `json:"items" binding:"required,gt=0,dive"`
}

type ItemResponseDTO struct {
	ID       uint     `json:"id"`
	OrderID  uint     `json:"order_id"`
	DishID   uint     `json:"dish_id"`
	Dish     dish.DTO `json:"dish"`
	Quantity int      `json:"quantity"`
}

type StatusDTO struct {
	ID    uint   `json:"id" binding:"required,min=1"`
	Title string `json:"title" binding:"required,min=2,max=30"`
}
