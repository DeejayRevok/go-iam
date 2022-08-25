package transformers

import (
	"go-uaa/src/domain/role"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/dto"
)

type UserToResponseTransformer struct {
	roleTransformer *RoleToResponseTransformer
}

func (transformer *UserToResponseTransformer) Transform(user *user.User) *dto.UserResponseDTO {
	responseDTO := dto.UserResponseDTO{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Roles:     transformer.transformRoles(user.Roles),
		Superuser: user.Superuser,
	}
	return &responseDTO
}

func (transformer *UserToResponseTransformer) transformRoles(roles []role.Role) []dto.RoleResponseDTO {
	var roleResponses []dto.RoleResponseDTO
	for _, role := range roles {
		roleResponses = append(roleResponses, *transformer.roleTransformer.Transform(&role))
	}
	return roleResponses
}

func NewUserToResponseTransformer(roleTransformer *RoleToResponseTransformer) *UserToResponseTransformer {
	transformer := UserToResponseTransformer{
		roleTransformer: roleTransformer,
	}
	return &transformer
}
