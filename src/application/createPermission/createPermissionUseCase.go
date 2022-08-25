package createPermission

import (
	"fmt"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/permission"

	"go.uber.org/zap"
)

type CreatePermissionUseCase struct {
	permissionRepository permission.PermissionRepository
	logger               *zap.Logger
}

func (useCase *CreatePermissionUseCase) Execute(request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*CreatePermissionRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(fmt.Sprintf("Starting creating permission with name %s", validatedRequest.Name))
	defer useCase.logger.Info(fmt.Sprintf("FinishedCreating permission with name %s", validatedRequest.Name))

	permission := permission.Permission{
		Name: validatedRequest.Name,
	}
	if err := useCase.permissionRepository.Save(permission); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (*CreatePermissionUseCase) RequiredPermissions() []string {
	return []string{permission.CreatePermissionPermission}
}

func NewCreatePermissionUseCase(permissionRepository permission.PermissionRepository, logger *zap.Logger) *CreatePermissionUseCase {
	useCase := CreatePermissionUseCase{
		permissionRepository: permissionRepository,
		logger:               logger,
	}
	return &useCase
}
