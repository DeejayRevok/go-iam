package dto

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoDTODeserializer struct{}

func (deserializer *EchoDTODeserializer) Deserialize(context echo.Context, destination interface{}) error {
	if err := context.Bind(destination); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error decoding request body")
	}
	if err := context.Validate(destination); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewEchoDTODeserializer() *EchoDTODeserializer {
	deserializer := EchoDTODeserializer{}
	return &deserializer
}
