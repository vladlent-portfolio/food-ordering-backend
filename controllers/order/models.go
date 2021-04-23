package order

import (
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/user"
	"time"
)

type Order struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	StatusID  uint
	Status    Status
	UserID    uint
	User      user.User `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Total     float64   `gorm:"check:total >= 0"`
	Items     []Item
}

type Item struct {
	ID       uint `gorm:"primaryKey"`
	OrderID  uint
	DishID   uint
	Dish     dish.Dish
	Quantity int `gorm:"check:quantity > 0"`
}

type Status struct {
	ID    uint   `gorm:"primaryKey"`
	Title string `gorm:"not null,unique,size:30"`
}

func (i Item) TableName() string {
	return "order_items"
}

func (s Status) TableName() string {
	return "order_status"
}
