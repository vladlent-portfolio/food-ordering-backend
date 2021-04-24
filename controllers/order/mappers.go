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
		Status:    o.Status,
		UserID:    o.UserID,
		User:      user.ToResponseDTO(o.User),
		Total:     o.Total,
		Items:     ToItemsResponseDTO(o.Items),
	}
}

func ToResponseDTOs(orders []Order) []ResponseDTO {
	dtos := make([]ResponseDTO, len(orders))

	for i, order := range orders {
		dtos[i] = ToResponseDTO(order)
	}

	return dtos
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

func ToItemsResponseDTO(items []Item) []ItemResponseDTO {
	dtos := make([]ItemResponseDTO, len(items))

	for i, item := range items {
		dtos[i] = ToItemResponseDTO(item)
	}

	return dtos
}
