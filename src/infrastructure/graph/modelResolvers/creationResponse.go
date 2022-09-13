package modelResolvers

type CreationResponse struct {
	success bool
}

func (response *CreationResponse) Success() *bool {
	return &response.success
}

func NewSuccessfulCreationResponse() *CreationResponse {
	return &CreationResponse{
		success: true,
	}
}

func NewFailedCreationResponse() *CreationResponse {
	return &CreationResponse{
		success: false,
	}
}
