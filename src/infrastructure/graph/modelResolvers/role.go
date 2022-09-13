package modelResolvers

import "go-uaa/src/infrastructure/dto"

type RoleResolver struct {
	role dto.RoleResponseDTO
}

func (resolver *RoleResolver) Name() *string {
	return &resolver.role.Name
}

func NewRoleResolver(role dto.RoleResponseDTO) *RoleResolver {
	return &RoleResolver{
		role: role,
	}
}
