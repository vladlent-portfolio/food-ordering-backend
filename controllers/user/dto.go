package user

import (
	"food_ordering_backend/common"
	"time"
)

type AuthDTO struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type ResponseDTO struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsAdmin   bool      `json:"is_admin"`
}

type DTOsWithPagination struct {
	Users      []ResponseDTO        `json:"users"`
	Pagination common.PaginationDTO `json:"pagination"`
}
