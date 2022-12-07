package modelResolvers

import (
	"go-iam/src/infrastructure/dto"
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

func NewUserResolver(user dto.UserResponseDTO) *UserResolver {
	return &UserResolver{
		user: user,
	}
}
