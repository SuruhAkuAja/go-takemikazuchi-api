package dto

import (
	"go-takemikazuchi-api/internal/user/dto"
	userAddressDto "go-takemikazuchi-api/internal/user_address/dto"
	workerDto "go-takemikazuchi-api/internal/worker/dto"
)

type JobResponseDto struct {
	ID           uint64                              `json:"id"`
	Title        string                              `json:"title"`
	Description  string                              `json:"description"`
	Price        float64                             `json:"price"`
	Status       string                              `json:"status"`
	CategoryName string                              `json:"category_name"`
	CreatedAt    string                              `json:"created_at"`
	UpdatedAt    string                              `json:"updated_at"`
	User         *dto.UserResponseDto                `json:"user"`
	Worker       *workerDto.WorkerResponseDto        `json:"worker"`
	UserAddress  *userAddressDto.UserAddressResponse `json:"user_address"`
}
