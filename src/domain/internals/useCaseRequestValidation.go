package internals

import "errors"

func ValidateUseCaseRequest[T any](request any) (T, *UseCaseResponse) {
	validatedRequest, requestOk := request.(T)
	if !requestOk {
		var zeroValue T
		errorResponse := ErrorUseCaseResponse(errors.New("malformed use case request"))
		return zeroValue, &errorResponse
	}
	return validatedRequest, nil
}
