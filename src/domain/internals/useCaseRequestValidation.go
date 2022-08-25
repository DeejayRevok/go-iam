package internals

import "errors"

func ValidateUseCaseRequest[T any](request any) (T, *UseCaseResponse) {
	validatedRequest, requestOk := request.(T)
	if requestOk == false {
		var zeroValue T
		errorResponse := ErrorUseCaseResponse(errors.New("Malformed use case request"))
		return zeroValue, &errorResponse
	}
	return validatedRequest, nil
}
