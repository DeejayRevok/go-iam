package controllers

import (
	"go-iam/src/infrastructure/dto"
	"go-iam/src/infrastructure/jwt"
	"go-iam/src/infrastructure/transformers"

	"github.com/labstack/echo/v4"
)

type GetJWTKeySetController struct {
	jwtKeySetBuilder *jwt.JWTKeySetBuilder
	dtoSerializer    *dto.EchoDTOSerializer
	errorTransformer *transformers.ErrorToEchoErrorTransformer
}

func (controller *GetJWTKeySetController) Handle(c echo.Context) error {
	keySet, err := controller.jwtKeySetBuilder.Build()
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	return controller.dtoSerializer.Serialize(c, keySet)
}

func NewGetJWTKeySetController(jwtKeySetBuilder *jwt.JWTKeySetBuilder, dtoSerializer *dto.EchoDTOSerializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *GetJWTKeySetController {
	controller := GetJWTKeySetController{
		jwtKeySetBuilder: jwtKeySetBuilder,
		dtoSerializer:    dtoSerializer,
		errorTransformer: errorTransformer,
	}
	return &controller
}
