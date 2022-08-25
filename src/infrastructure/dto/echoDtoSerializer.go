package dto

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoDTOSerializer struct{}

func (serializer *EchoDTOSerializer) Serialize(context echo.Context, source interface{}) error {
	return context.JSON(http.StatusOK, source)
}

func NewEchoDTOSerializer() *EchoDTOSerializer {
	serializer := EchoDTOSerializer{}
	return &serializer
}
