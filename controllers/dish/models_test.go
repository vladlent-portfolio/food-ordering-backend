package dish

import (
	"food_ordering_backend/controllers/category"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestCategories = []category.Category{
	{ID: 1, Title: "Salads", Removable: true},
	{ID: 2, Title: "Burgers", Removable: true},
	{ID: 3, Title: "Pizza", Removable: true},
	{ID: 4, Title: "Drinks", Removable: true},
}
var TestDishes = []Dish{
	{ID: 1, Title: "Fresh and Healthy Salad", Price: 2.65, CategoryID: 1, Category: TestCategories[0]},
	{ID: 2, Title: "Crunchy Cashew Salad", Price: 3.22, CategoryID: 1, Category: TestCategories[0]},
	{ID: 3, Title: "Hamburger", Price: 1.99, CategoryID: 2, Category: TestCategories[1]},
	{ID: 4, Title: "Cheeseburger", Price: 2.28, CategoryID: 2, Category: TestCategories[1]},
	{ID: 5, Title: "Margherita", Price: 4.20, CategoryID: 3, Category: TestCategories[2]},
	{ID: 6, Title: "4 Cheese", Price: 4.69, CategoryID: 3, Category: TestCategories[2]},
	{ID: 7, Title: "Pepsi 2L", Price: 1.50, CategoryID: 4, Category: TestCategories[3]},
	{ID: 8, Title: "Orange Juice 2L", Price: 2, CategoryID: 4, Category: TestCategories[3]},
}

func TestDishes_Find(t *testing.T) {
	dishes := Dishes(TestDishes)
	t.Run("should return dish and true", func(t *testing.T) {
		it := assert.New(t)
		tests := []struct {
			lookupFn func(d Dish, index int) bool
			expected Dish
		}{
			{
				lookupFn: func(d Dish, index int) bool {
					return d.ID == dishes[0].ID
				},
				expected: dishes[0],
			},
			{
				lookupFn: func(d Dish, index int) bool {
					return d.Title == dishes[1].Title
				},
				expected: dishes[1],
			},
			{
				lookupFn: func(d Dish, index int) bool {
					return d.Price == dishes[2].Price
				},
				expected: dishes[2],
			},
			{
				lookupFn: func(d Dish, index int) bool {
					return d.CategoryID == dishes[3].CategoryID
				},
				expected: dishes[2],
			},
		}

		for _, tc := range tests {
			d, ok := dishes.Find(tc.lookupFn)
			it.Equal(tc.expected, d)
			it.True(ok)
		}
	})

	t.Run("should return empty dish and false", func(t *testing.T) {
		it := assert.New(t)
		fns := []func(d Dish, index int) bool{
			func(d Dish, index int) bool {
				return d.ID == 32
			},
			func(d Dish, index int) bool {
				return d.Title == "Not exist"
			},
			func(d Dish, index int) bool {
				return d.CategoryID == 56
			},
			func(d Dish, index int) bool {
				return d.Price == 56
			},
		}
		for _, lookup := range fns {
			d, ok := dishes.Find(lookup)
			it.Zero(d)
			it.False(ok)
		}
	})
}
