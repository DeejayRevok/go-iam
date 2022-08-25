package internals

type UseCaseResponse struct {
	Content interface{}
	Err     error
}

func ErrorUseCaseResponse(err error) UseCaseResponse {
	return UseCaseResponse{
		Err: err,
	}
}

func EmptyUseCaseResponse() UseCaseResponse {
	return UseCaseResponse{
		Content: nil,
		Err:     nil,
	}
}
