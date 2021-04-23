package order

import (
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/user"
	"time"
)

type Status int

const (
	StatusCreated    Status = 0
	StatusInProgress Status = 1
	StatusDone       Status = 2
	StatusCanceled   Status = 3
)

type Order struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    Status `gorm:"type:smallint"`
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

func (i Item) TableName() string {
	return "order_items"
}
