package dish

import (
	"food_ordering_backend/controllers/category"
	"gorm.io/gorm"
)

func ToModel(dto DTO) Dish {
	return Dish{
		Model:      gorm.Model{ID: dto.ID},
		Title:      dto.Title,
		CategoryID: dto.CategoryID,
		Price:      dto.Price,
		Category:   category.ToModel(dto.Category),
	}
}

func ToDTO(d Dish) DTO {
	return DTO{
		ID:         d.ID,
		Title:      d.Title,
		Price:      d.Price,
		CategoryID: d.CategoryID,
		Category:   category.ToDTO(d.Category),
	}
}

func ToDTOs(dishes []Dish) []DTO {
	dtos := make([]DTO, len(dishes))

	for i, d := range dishes {
		dtos[i] = ToDTO(d)
	}

	return dtos
}
