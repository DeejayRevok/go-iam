package controllers

import (
	"go-iam/src/application/getApplicationHealth"
	"go-iam/src/domain/internals"
	"go-iam/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetStatusController struct {
	getApplicationHealthUseCase *getApplicationHealth.GetApplicationHealthUseCase
	useCaseExecutor             *internals.UseCaseExecutor
	errorTransformer            *transformers.ErrorToEchoErrorTransformer
}

func (controller *GetStatusController) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.getApplicationHealthUseCase, struct{}{}, nil)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusOK)
}

func NewGetStatusController(getApplicationHealthUseCase *getApplicationHealth.GetApplicationHealthUseCase, useCaseExecutor *internals.UseCaseExecutor, errorTransformer *transformers.ErrorToEchoErrorTransformer) *GetStatusController {
	controller := GetStatusController{
		getApplicationHealthUseCase: getApplicationHealthUseCase,
		useCaseExecutor:             useCaseExecutor,
		errorTransformer:            errorTransformer,
	}
	return &controller
}
