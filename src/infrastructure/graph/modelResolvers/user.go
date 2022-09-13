package modelResolvers

import (
	"go-uaa/src/infrastructure/dto"
)

type UserResolver struct {
	user dto.UserResponseDTO
}

func (resolver *UserResolver) Username() *string {
	return &resolver.user.Username
}

func (resolver *UserResolver) Email() *string {
	return &resolver.user.Email
}

func (resolver *UserResolver) Roles() *[]*RoleResolver {
	roles := make([]*RoleResolver, 0)
	for _, role := range resolver.user.Roles {
		roles = append(roles, NewRoleResolver(role))
	}
	return &roles
}

func NewUserResolver(user dto.UserResponseDTO) *UserResolver {
	return &UserResolver{
		user: user,
	}
}
