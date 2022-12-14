package controllers

import (
	"go-iam/src/application/requestPasswordReset"
	"go-iam/src/domain/internals"
	"go-iam/src/infrastructure/dto"
	"go-iam/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestResetPasswordController struct {
	requestResetPasswordUseCase *requestPasswordReset.RequestPasswordResetUseCase
	useCaseExecutor             *internals.UseCaseExecutor
	dtoDeserializer             *dto.EchoDTODeserializer
	errorTransformer            *transformers.ErrorToEchoErrorTransformer
}

func (controller *RequestResetPasswordController) Handle(c echo.Context) error {
	var requestResetDTO dto.RequestResetPasswordDTO
	err := controller.dtoDeserializer.Deserialize(c, &requestResetDTO)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	ctx := c.Request().Context()
	useCaseRequest := requestPasswordReset.RequestPasswordResetRequest{
		Email: requestResetDTO.Email,
	}
	if useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.requestResetPasswordUseCase, &useCaseRequest, nil); useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}

	return c.NoContent(http.StatusCreated)
}

func NewRequestResetPasswordController(useCase *requestPasswordReset.RequestPasswordResetUseCase, useCaseExecutor *internals.UseCaseExecutor, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *RequestResetPasswordController {
	return &RequestResetPasswordController{
		requestResetPasswordUseCase: useCase,
		useCaseExecutor:             useCaseExecutor,
		dtoDeserializer:             dtoDeserializer,
		errorTransformer:            errorTransformer,
	}
}
