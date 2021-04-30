package testutils

import (
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

var TestCategories = []category.Category{
	{ID: 1, Title: "Salads", Removable: true, Image: strPointer("1.png")},
	{ID: 2, Title: "Burgers", Removable: true, Image: strPointer("2.png")},
	{ID: 3, Title: "Pizza", Removable: true, Image: strPointer("3.png")},
	{ID: 4, Title: "Drinks", Removable: true, Image: strPointer("4.png")},
}
var TestDishes = []dish.Dish{
	{ID: 1, Title: "Fresh and Healthy Salad", Price: 2.65, CategoryID: 1, Category: TestCategories[0]},
	{ID: 2, Title: "Crunchy Cashew Salad", Price: 3.22, CategoryID: 1, Category: TestCategories[0]},
	{ID: 3, Title: "Hamburger", Price: 1.99, CategoryID: 2, Category: TestCategories[1]},
	{ID: 4, Title: "Cheeseburger", Price: 2.28, CategoryID: 2, Category: TestCategories[1]},
	{ID: 5, Title: "Margherita", Price: 4.20, CategoryID: 3, Category: TestCategories[2]},
	{ID: 6, Title: "4 Cheese", Price: 4.69, CategoryID: 3, Category: TestCategories[2]},
	{ID: 7, Title: "Pepsi 2L", Price: 1.50, CategoryID: 4, Category: TestCategories[3]},
	{ID: 8, Title: "Orange Juice 2L", Price: 2, CategoryID: 4, Category: TestCategories[3]},
}
var TestCategoriesJSON = `[{"id":1,"title":"Salads","removable":true},{"id":2,"title":"Burgers","removable":true},{"id":3,"title":"Pizza","removable":true},{"id":4,"title":"Drinks","removable":true}]`
var TestDishesJSON = `[{"id":1,"title":"Fresh and Healthy Salad","price":2.65,"category_id":1,"category":{"id":1,"title":"Salads","removable":true}},{"id":2,"title":"Crunchy Cashew Salad","price":3.22,"category_id":1,"category":{"id":1,"title":"Salads","removable":true}},{"id":3,"title":"Hamburger","price":1.99,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}},{"id":4,"title":"Cheeseburger","price":2.28,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}},{"id":5,"title":"Margherita","price":4.2,"category_id":3,"category":{"id":3,"title":"Pizza","removable":true}},{"id":6,"title":"4 Cheese","price":4.69,"category_id":3,"category":{"id":3,"title":"Pizza","removable":true}},{"id":7,"title":"Pepsi 2L","price":1.5,"category_id":4,"category":{"id":4,"title":"Drinks","removable":true}},{"id":8,"title":"Orange Juice 2L","price":2,"category_id":4,"category":{"id":4,"title":"Drinks","removable":true}}]`

func SetupDishesAndCategories(t *testing.T) {
	req := require.New(t)
	cleanup := func() {
		req.NoError(db.Exec("TRUNCATE categories CASCADE;").Error)
	}
	cleanup()
	t.Cleanup(cleanup)

	req.NoError(db.Create(TestDishes).Error)
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

func strPointer(str string) *string {
	return &str
}
