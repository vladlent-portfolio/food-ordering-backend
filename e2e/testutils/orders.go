package testutils

import (
	"food_ordering_backend/controllers/order"
	"food_ordering_backend/controllers/user"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"testing"
	"time"
)

type orderGenerator struct {
	orderID     uint
	orderItemID uint
}

var gen = &orderGenerator{1, 1}

var TestOrders = []order.Order{
	gen.createTestOrder(order.StatusDone, TestUsers[0]),
	gen.createTestOrder(order.StatusInProgress, TestUsers[0]),
	gen.createTestOrder(order.StatusCanceled, TestUsers[1]),
	gen.createTestOrder(order.StatusCreated, TestUsers[2]),
	gen.createTestOrder(order.StatusCreated, TestUsers[0]),
}

var TestOrderItems = []order.Item{
	gen.createOrderItem(1, 1, 1),
	gen.createOrderItem(1, 4, 2),
	gen.createOrderItem(1, 7, 1),
	gen.createOrderItem(2, 5, 3),
	gen.createOrderItem(2, 8, 2),
	gen.createOrderItem(3, 2, 4),
	gen.createOrderItem(3, 3, 4),
	gen.createOrderItem(3, 7, 2),
	gen.createOrderItem(3, 8, 2),
	gen.createOrderItem(4, 6, 1),
	gen.createOrderItem(4, 7, 2),
	gen.createOrderItem(5, 5, 3),
}

// SetupOrdersDB populates db with TestOrders and TestOrderItems.
// It will call SetupUsersDB and SetupDishesAndCategories before creating orders
// and will reset all tables after test run.
func SetupOrdersDB(t *testing.T) {
	req := require.New(t)
	cleanup := func() {
		req.NoError(db.Exec("TRUNCATE orders, order_items;").Error)
	}
	t.Cleanup(cleanup)
	cleanup()

	SetupUsersDB(t)
	SetupDishesAndCategories(t)

	rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(TestOrders), func(i, j int) {
		TestOrders[i], TestOrders[j] = TestOrders[j], TestOrders[i]
	})

	req.NoError(db.Create(&TestOrders).Error)
	//req.NoError(db.Create(&TestOrderItems).Error)
}

func (g *orderGenerator) createTestOrder(status order.Status, u user.User) order.Order {
	o := order.Order{
		ID:     g.orderID,
		Status: status,
		UserID: u.ID,
		User:   u,
		Items:  findOrderItemsByOrderID(g.orderID),
	}
	g.orderID++

	o.Total = order.CalcTotal(o.Items)
	return o
}

func (g *orderGenerator) createOrderItem(orderID, testDishID uint, quantity int) order.Item {
	d := FindTestDishByID(testDishID)
	i := order.Item{
		ID:       g.orderItemID,
		OrderID:  orderID,
		Dish:     d,
		DishID:   d.ID,
		Quantity: quantity,
	}

	g.orderItemID++
	return i
}

func FindTestOrderByID(id uint) order.Order {
	for _, o := range TestOrders {
		if o.ID == id {
			return o
		}
	}
	log.Panicf("cannot find TestOrder with id %d\n", id)
	return order.Order{}
}

func findOrderItemsByOrderID(id uint) []order.Item {
	var items []order.Item
	for _, item := range TestOrderItems {
		if item.OrderID == id {
			items = append(items, item)
		}
	}
	return items
}
