package controllers

import (
	"go-iam/src/application/resetPassword"
	"go-iam/src/domain/internals"
	"go-iam/src/infrastructure/dto"
	"go-iam/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResetPasswordController struct {
	resetPasswordUseCase *resetPassword.ResetPasswordUseCase
	useCaseExecutor      *internals.UseCaseExecutor
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
	if useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.resetPasswordUseCase, &useCaseRequest, nil); useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusOK)
}

func NewResetPasswordController(useCase *resetPassword.ResetPasswordUseCase, useCaseExecutor *internals.UseCaseExecutor, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *ResetPasswordController {
	return &ResetPasswordController{
		resetPasswordUseCase: useCase,
		useCaseExecutor:      useCaseExecutor,
		dtoDeserializer:      dtoDeserializer,
		errorTransformer:     errorTransformer,
	}
}
