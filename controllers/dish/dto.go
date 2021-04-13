package dish

import "food_ordering_backend/controllers/category"

type DTO struct {
	ID         uint         `json:"id,omitempty"`
	Title      string       `json:"title"`
	Price      float64      `json:"price"`
	CategoryID uint         `json:"category_id"`
	Category   category.DTO `json:"category,omitempty"`
}
