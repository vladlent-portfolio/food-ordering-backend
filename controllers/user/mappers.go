package user

import "golang.org/x/crypto/bcrypt"

func CreateFromDTO(dto AuthDTO) (User, error) {
	var user User
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.Email = dto.Email
	user.PasswordHash = hashedPass

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
