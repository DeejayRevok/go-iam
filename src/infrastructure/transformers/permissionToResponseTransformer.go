package transformers

import (
	"go-uaa/src/domain/permission"
	"go-uaa/src/infrastructure/dto"
)

type PermissionToResponseTransformer struct{}

func (transformer *PermissionToResponseTransformer) Transform(permission *permission.Permission) *dto.PermissionResponseDTO {
	permissionResponse := dto.PermissionResponseDTO{
		Name: permission.Name,
	}
	return &permissionResponse
}

func NewPermissionToResponseTransformer() *PermissionToResponseTransformer {
	transformer := PermissionToResponseTransformer{}
	return &transformer
}
