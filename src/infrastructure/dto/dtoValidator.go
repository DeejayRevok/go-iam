package dto

import (
	"github.com/go-playground/validator"
)

type DTOValidator struct {
	validator *validator.Validate
}

func (v *DTOValidator) Validate(dto interface{}) error {
	if err := v.validator.Struct(dto); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return DTOValidationError{
			validation: validationErrors[0].Tag(),
			field:      validationErrors[0].Field(),
		}
	}
	return nil
}

func NewDTOValidator() *DTOValidator {
	dtoValidator := DTOValidator{
		validator: validator.New(),
	}
	return &dtoValidator
}
