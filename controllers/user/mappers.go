package user

func ToResponseDTO(u User) ResponseDTO {
	return ResponseDTO{
		ID:    u.ID,
		Email: u.Email,
	}
}
