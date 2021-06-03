package order

import (
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/user"
	"math"
	"time"
)

type Status int

const (
	StatusCreated    Status = 0
	StatusInProgress Status = 1
	StatusDone       Status = 2
	StatusCanceled   Status = 3
)

var Statuses = []Status{StatusCreated, StatusInProgress, StatusDone, StatusCanceled}

type Items []Item

type Order struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    Status `gorm:"type:smallint;check:status IN (0,1,2,3)"`
	UserID    uint
	User      user.User `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Total     float64   `gorm:"check:total >= 0"`
	Items     []Item    `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type Item struct {
	ID       uint `gorm:"primaryKey"`
	OrderID  uint
	DishID   uint
	Dish     dish.Dish `gorm:"constraint:OnUpdate:CASCADE"`
	Quantity int       `gorm:"type:int;check:quantity > 0"`
}

func (i Item) TableName() string {
	return "order_items"
}

// Cost calculates the total cost of order item.
func (i *Item) Cost() float64 {
	res := i.Dish.Price * float64(i.Quantity)
	// Dealing with precision problems
	return math.Ceil(res*100) / 100
}

// IDs is a convenience method to extract all ids from Items.
func (items Items) IDs() []uint {
	ids := make([]uint, len(items))

	for i, item := range items {
		ids[i] = item.ID
	}

	return ids
}

// CalcTotal calculates the total of cost of all items in a provided slice.
func CalcTotal(items []Item) float64 {
	var total float64

	for _, item := range items {
		total += item.Cost()
	}

	return total
}

// IsValidStatus checks whether provided status is a valid Status.
// Useful to validate input that comes from external sources, e.g as
// a query parameter.
func IsValidStatus(status int) bool {
	for _, s := range Statuses {
		if status == int(s) {
			return true
		}
	}

	return false
}
