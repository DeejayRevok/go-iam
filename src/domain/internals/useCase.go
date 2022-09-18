package internals

import "context"

type UseCase interface {
	Execute(ctx context.Context, request interface{}) UseCaseResponse
	RequiredPermissions() []string
}
