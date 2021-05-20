package testutils

import (
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

// TODO: Shuffle categories and dishes

var TestCategories = []category.Category{
	{ID: 3, Title: "Pizza", Removable: true, Image: strPointer("3.png")},
	{ID: 1, Title: "Salads", Removable: true, Image: strPointer("1.png")},
	{ID: 4, Title: "Drinks", Removable: true, Image: strPointer("4.png")},
	{ID: 2, Title: "Burgers", Removable: true, Image: strPointer("2.png")},
}
var TestDishes = []dish.Dish{
	{ID: 1, Title: "Fresh and Healthy Salad", Price: 2.65, Image: strPointer("1.png"), CategoryID: 1, Category: FindTestCategoryByID(1)},
	{ID: 2, Title: "Crunchy Cashew Salad", Price: 3.22, Image: strPointer("2.png"), CategoryID: 1, Category: FindTestCategoryByID(1)},
	{ID: 3, Title: "Hamburger", Price: 1.99, Image: strPointer("3.png"), CategoryID: 2, Category: FindTestCategoryByID(2)},
	{ID: 4, Title: "Cheeseburger", Price: 2.28, Image: strPointer("4.png"), CategoryID: 2, Category: FindTestCategoryByID(2)},
	{ID: 5, Title: "Margherita", Price: 4.20, Image: strPointer("5.png"), CategoryID: 3, Category: FindTestCategoryByID(3)},
	{ID: 6, Title: "4 Cheese", Price: 4.69, Image: strPointer("6.png"), CategoryID: 3, Category: FindTestCategoryByID(3)},
	{ID: 7, Title: "Pepsi 2L", Price: 1.50, Image: strPointer("7.png"), CategoryID: 4, Category: FindTestCategoryByID(4)},
	{ID: 8, Title: "Orange Juice 2L", Price: 2, Image: strPointer("8.png"), CategoryID: 4, Category: FindTestCategoryByID(4)},
}

func SetupDishesAndCategories(t *testing.T) {
	req := require.New(t)
	cleanup := func() {
		req.NoError(db.Exec("TRUNCATE categories CASCADE;").Error)
	}
	cleanup()
	t.Cleanup(cleanup)

	// Inserting categories separately, to preserve their order
	req.NoError(db.Create(TestCategories).Error)
	req.NoError(db.Omit("Category").Create(TestDishes).Error)
}

func FindTestDishByID(id uint) dish.Dish {
	for _, testDish := range TestDishes {
		if testDish.ID == id {
			return testDish
		}
	}
	log.Panicf("cannot find TestDish with id %d\n", id)
	return dish.Dish{}
}

func FindTestCategoryByID(id uint) category.Category {
	for _, c := range TestCategories {
		if c.ID == id {
			return c
		}
	}
	log.Panicf("cannot find TestCategory with id %d\n", id)
	return category.Category{}
}

func strPointer(str string) *string {
	return &str
}
