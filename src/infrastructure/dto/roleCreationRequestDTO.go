package dto

type RoleCreationRequestDTO struct {
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}
