package dto

type UserResponseDto struct {
	Id        uint64 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at" mapstructure:"-"`
	UpdatedAt string `json:"updated_at" mapstructure:"-"`
}
