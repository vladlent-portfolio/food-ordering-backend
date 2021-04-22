package order

import (
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
	User      user.User
	Item      Item
	Total     float64 `gorm:"check:total >= 0"`
}

type Item struct {
	ID uint `gorm:"primaryKey"`
}

type Status struct {
	ID    uint   `gorm:"primaryKey"`
	Title string `gorm:"not null,unique"`
}

func (o Order) TableName() string {
	return "order.orders"
}

func (s Status) TableName() string {
	return "order.status"
}
