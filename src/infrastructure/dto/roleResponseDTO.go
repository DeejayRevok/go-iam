package dto

type RoleResponseDTO struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Permissions []PermissionResponseDTO `json:"permissions"`
}
