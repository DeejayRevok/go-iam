package internals

type UseCase interface {
	Execute(request interface{}) UseCaseResponse
	RequiredPermissions() []string
}
