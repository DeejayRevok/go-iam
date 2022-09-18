package createRole

import (
	"context"
	"fmt"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/permission"
	"go-uaa/src/domain/role"

	"github.com/google/uuid"
)

type CreateRoleUseCase struct {
	roleRepository       role.RoleRepository
	permissionRepository permission.PermissionRepository
	logger               internals.Logger
}

func (useCase *CreateRoleUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*CreateRoleRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting creation of role %s", validatedRequest.Name))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished creation of role %s", validatedRequest.Name))

	permissions, err := useCase.findPermissions(ctx, validatedRequest.Permissions)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	role := role.Role{
		ID:          uuid.New(),
		Name:        validatedRequest.Name,
		Permissions: permissions,
	}
	if err = useCase.roleRepository.Save(ctx, role); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (useCase *CreateRoleUseCase) findPermissions(ctx context.Context, permissionNames []string) ([]permission.Permission, error) {
	permissions, err := useCase.permissionRepository.FindByNames(ctx, permissionNames)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (*CreateRoleUseCase) RequiredPermissions() []string {
	return []string{role.CreateRolePermission}
}

func NewCreateRoleUseCase(roleRepository role.RoleRepository, permissionRepository permission.PermissionRepository, logger internals.Logger) *CreateRoleUseCase {
	useCase := CreateRoleUseCase{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		logger:               logger,
	}
	return &useCase
}
