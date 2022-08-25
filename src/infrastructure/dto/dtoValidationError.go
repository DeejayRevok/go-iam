package dto

import "fmt"

type DTOValidationError struct {
	validation string
	field      string
}

func (err DTOValidationError) Error() string {
	return fmt.Sprintf("Validation '%s' failed for field %s", err.validation, err.field)
}
