package dish

import (
	"food_ordering_backend/controllers/category"
	"path"
)

func ToModel(dto DTO) Dish {
	image := dto.Image

	if image != nil {
		name := path.Base(*image)
		image = &name
	}

	return Dish{
		ID:         dto.ID,
		Title:      dto.Title,
		CategoryID: dto.CategoryID,
		Price:      dto.Price,
		Image:      image,
		Category:   category.ToModel(dto.Category),
	}
}

func ToDTO(d Dish) DTO {
	image := d.Image

	if image != nil {
		uri := PathToImg(*image)
		image = &uri
	}

	return DTO{
		ID:         d.ID,
		Title:      d.Title,
		Price:      d.Price,
		CategoryID: d.CategoryID,
		Image:      image,
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
