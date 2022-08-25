package dto

type UserResponseDTO struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Roles     []RoleResponseDTO `json:"roles"`
	Superuser bool              `json:"superuser"`
}
