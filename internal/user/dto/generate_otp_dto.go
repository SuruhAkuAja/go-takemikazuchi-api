package dto

type GenerateOtpDto struct {
	UserId uint64 `json:"user_id"`
	Email  string `json:"email" validate:"required,email"`
	Role   string `json:"role" validate:"required,role"`
}
