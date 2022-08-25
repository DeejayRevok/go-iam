package transformers

import (
	"go-uaa/src/domain/permission"
	"go-uaa/src/domain/role"
	"go-uaa/src/infrastructure/dto"
)

type RoleToResponseTransformer struct {
	permissionTransformer *PermissionToResponseTransformer
}

func (transformer *RoleToResponseTransformer) Transform(role *role.Role) *dto.RoleResponseDTO {
	roleResponse := dto.RoleResponseDTO{
		ID:          role.ID.String(),
		Name:        role.Name,
		Permissions: transformer.transformPermissions(role.Permissions),
	}
	return &roleResponse
}

func (transformer *RoleToResponseTransformer) transformPermissions(permissions []permission.Permission) []dto.PermissionResponseDTO {
	var permissionResponses []dto.PermissionResponseDTO
	for _, permission := range permissions {
		permissionResponses = append(permissionResponses, *transformer.permissionTransformer.Transform(&permission))
	}
	return permissionResponses
}

func NewRoleToResponseTransformer(permissionTransformer *PermissionToResponseTransformer) *RoleToResponseTransformer {
	transformer := RoleToResponseTransformer{
		permissionTransformer: permissionTransformer,
	}
	return &transformer
}
