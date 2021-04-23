package order

import (
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/user"
)

func ToResponseDTO(o Order) ResponseDTO {
	return ResponseDTO{
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		StatusID:  o.StatusID,
		Status:    ToStatusDTO(o.Status),
		UserID:    o.UserID,
		User:      user.ToResponseDTO(o.User),
		Total:     o.Total,
		Items:     ToItemsDTO(o.Items),
	}
}

// TODO: Change all ToModel() functions to FromDTO()

//func CreateFromRequestDTO(dto RequestDTO) Order {
//	var o Order
//	o.Items = dto.Items
//}

func ToItemResponseDTO(i Item) ItemResponseDTO {
	return ItemResponseDTO{
		ID:       i.ID,
		OrderID:  i.OrderID,
		DishID:   i.DishID,
		Dish:     dish.ToDTO(i.Dish),
		Quantity: i.Quantity,
	}
}

func ToItemsDTO(items []Item) []ItemResponseDTO {
	dtos := make([]ItemResponseDTO, len(items))

	for i, item := range items {
		dtos[i] = ToItemResponseDTO(item)
	}

	return dtos
}

func ToStatusDTO(s Status) StatusDTO {
	return StatusDTO{
		ID:    s.ID,
		Title: s.Title,
	}
}
