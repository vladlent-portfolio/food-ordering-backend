package category

func ToCategory(dto DTO) Category {
	return Category{Title: dto.Title, Removable: dto.Removable}
}

func ToDTO(c Category) DTO {
	return DTO{ID: c.ID, Title: c.Title, Removable: c.Removable}
}

func ToCategoriesDTOs(categories []Category) []DTO {
	dtos := make([]DTO, len(categories))

	for i, c := range categories {
		dtos[i] = ToDTO(c)
	}

	return dtos
}
