package category

func ToModel(dto DTO) Category {
	return Category{
		ID:        dto.ID,
		Title:     dto.Title,
		Removable: dto.Removable,
		Image:     dto.Image,
	}
}

func ToDTO(c Category) DTO {
	return DTO{
		ID:        c.ID,
		Title:     c.Title,
		Removable: c.Removable,
		Image:     c.Image,
	}
}

func ToDTOs(categories []Category) []DTO {
	dtos := make([]DTO, len(categories))

	for i, c := range categories {
		dtos[i] = ToDTO(c)
	}

	return dtos
}
