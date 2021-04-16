package user

func CreateFromDTO(dto AuthDTO) (User, error) {
	var user User
	err := user.SetPassword(dto.Password)

	if err != nil {
		return user, err
	}

	user.Email = dto.Email

	return user, nil
}

func ToResponseDTO(u User) ResponseDTO {
	return ResponseDTO{
		ID:    u.ID,
		Email: u.Email,
	}
}

func ToResponseDTOs(users []User) []ResponseDTO {
	dtos := make([]ResponseDTO, len(users))

	for i, u := range users {
		dtos[i] = ToResponseDTO(u)
	}

	return dtos
}
