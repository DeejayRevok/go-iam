package dto

type PermissionCreationRequestDTO struct {
	Name string `json:"name" validate:"required"`
}
