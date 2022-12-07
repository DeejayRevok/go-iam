package transformers

import (
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/dto"
)

type UserToResponseTransformer struct{}

func (transformer *UserToResponseTransformer) Transform(user *user.User) *dto.UserResponseDTO {
	responseDTO := dto.UserResponseDTO{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Superuser: user.Superuser,
	}
	return &responseDTO
}

func NewUserToResponseTransformer() *UserToResponseTransformer {
	transformer := UserToResponseTransformer{}
	return &transformer
}
