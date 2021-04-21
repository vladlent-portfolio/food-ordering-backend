package user

func CreateFromDTO(dto AuthDTO) (User, error) {
	var user User

	user.Email = dto.Email
	user.SetPassword(dto.Password)

	return user, nil
}

func ToResponseDTO(u User) ResponseDTO {
	return ResponseDTO{
		ID:      u.ID,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}

func ToResponseDTOs(users []User) []ResponseDTO {
	dtos := make([]ResponseDTO, len(users))

	for i, u := range users {
		dtos[i] = ToResponseDTO(u)
	}

	return dtos
}
