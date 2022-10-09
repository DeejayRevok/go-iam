package controllers

import (
	"go-uaa/src/application/resetPassword"
	"go-uaa/src/domain/internals"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResetPasswordController struct {
	resetPasswordUseCase *resetPassword.ResetPasswordUseCase
	useCaseExecutor      *internals.AuthorizedUseCaseExecutor
	dtoDeserializer      *dto.EchoDTODeserializer
	errorTransformer     *transformers.ErrorToEchoErrorTransformer
}

func (controller *ResetPasswordController) Handle(c echo.Context) error {
	var resetDTO dto.ResetPasswordDTO
	if err := controller.dtoDeserializer.Deserialize(c, &resetDTO); err != nil {
		return controller.errorTransformer.Transform(err)
	}

	ctx := c.Request().Context()
	useCaseRequest := resetPassword.ResetPasswordRequest{
		UserEmail:   resetDTO.UserEmail,
		ResetToken:  resetDTO.ResetToken,
		NewPassword: resetDTO.NewPassword,
	}
	if useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.resetPasswordUseCase, &useCaseRequest, nil, nil); useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusOK)
}

func NewResetPasswordController(useCase *resetPassword.ResetPasswordUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *ResetPasswordController {
	return &ResetPasswordController{
		resetPasswordUseCase: useCase,
		useCaseExecutor:      useCaseExecutor,
		dtoDeserializer:      dtoDeserializer,
		errorTransformer:     errorTransformer,
	}
}
