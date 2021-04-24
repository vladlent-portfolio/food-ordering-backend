package testutils

import (
	"food_ordering_backend/common"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/order"
)

var TestOrders []order.Order

func init() {

}

func populateTestOrders() {

}

func randomOrder(orderID uint) order.Order {
	u := TestUsers[common.RandomInt(len(TestUsers))]

	return order.Order{
		ID:     orderID,
		Status: order.StatusCreated,
		UserID: u.ID,
		User:   u,
	}
}

func randomItem(orderID uint) order.Item {
	return order.Item{
		ID:       0,
		OrderID:  orderID,
		DishID:   0,
		Dish:     dish.Dish{},
		Quantity: 0,
	}
}
