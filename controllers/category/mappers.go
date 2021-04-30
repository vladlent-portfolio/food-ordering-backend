package category

import (
	"path"
)

func ToModel(dto DTO) Category {
	image := dto.Image

	if image != nil {
		name := path.Base(*image)
		image = &name
	}

	return Category{
		ID:        dto.ID,
		Title:     dto.Title,
		Removable: dto.Removable,
		Image:     image,
	}
}

func ToDTO(c Category) DTO {
	image := c.Image

	if image != nil {
		uri := PathToImg(*image)
		image = &uri
	}

	return DTO{
		ID:        c.ID,
		Title:     c.Title,
		Removable: c.Removable,
		Image:     image,
	}
}

func ToDTOs(categories []Category) []DTO {
	dtos := make([]DTO, len(categories))

	for i, c := range categories {
		dtos[i] = ToDTO(c)
	}

	return dtos
}
