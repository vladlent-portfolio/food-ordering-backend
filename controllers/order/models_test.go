package order

import (
	"food_ordering_backend/controllers/dish"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_Cost(t *testing.T) {
	t.Run("should correctly calculate the cost of order item", func(t *testing.T) {
		tests := []struct {
			item     Item
			expected float64
		}{
			{item: Item{Dish: dish.Dish{Price: 3.22}, Quantity: 3}, expected: 9.66},
			{item: Item{Dish: dish.Dish{Price: 2.28}, Quantity: 10}, expected: 22.8},
			{item: Item{Dish: dish.Dish{Price: 4.20}, Quantity: 4}, expected: 16.8},
			{item: Item{Dish: dish.Dish{Price: 0.3}, Quantity: 3}, expected: 0.9},
			{item: Item{Dish: dish.Dish{Price: 1.1}, Quantity: 3}, expected: 3.3},
		}

		for _, tc := range tests {
			assert.Equal(t, tc.expected, tc.item.Cost())
		}

	})
}

func TestCalcTotal(t *testing.T) {
	t.Run("should correctly calculate total cost", func(t *testing.T) {
		tests := []struct {
			items    []Item
			expected float64
		}{
			{
				items: []Item{
					{Dish: dish.Dish{Price: 3.22}, Quantity: 3},
					{Dish: dish.Dish{Price: 0.3}, Quantity: 3},
				},
				expected: 10.56,
			},
			{
				items: []Item{
					{Dish: dish.Dish{Price: 2.28}, Quantity: 10},
					{Dish: dish.Dish{Price: 4.20}, Quantity: 4},
					{Dish: dish.Dish{Price: 1.1}, Quantity: 3},
				},
				expected: 42.9,
			},
			{
				items: []Item{
					{Dish: dish.Dish{Price: 3.22}, Quantity: 3},
					{Dish: dish.Dish{Price: 2.28}, Quantity: 10},
					{Dish: dish.Dish{Price: 4.20}, Quantity: 4},
					{Dish: dish.Dish{Price: 0.3}, Quantity: 3},
					{Dish: dish.Dish{Price: 1.1}, Quantity: 3},
				},
				expected: 53.46,
			},
		}

		for _, tc := range tests {
			assert.Equal(t, tc.expected, CalcTotal(tc.items))
		}
	})
}

func TestItems_IDs(t *testing.T) {
	items := Items{{ID: 23}, {ID: 456}, {ID: 10}}

	t.Run("should return a slice of ids in original order", func(t *testing.T) {
		assert.Equal(t, []uint{23, 456, 10}, items.IDs())
	})
}
