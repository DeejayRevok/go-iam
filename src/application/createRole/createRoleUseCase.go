package createRole

import (
	"fmt"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/permission"
	"go-uaa/src/domain/role"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateRoleUseCase struct {
	roleRepository       role.RoleRepository
	permissionRepository permission.PermissionRepository
	logger               *zap.Logger
}

func (useCase *CreateRoleUseCase) Execute(request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*CreateRoleRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(fmt.Sprintf("Starting creation of role %s", validatedRequest.Name))
	defer useCase.logger.Info(fmt.Sprintf("Finished creation of role %s", validatedRequest.Name))

	permissions, err := useCase.findPermissions(validatedRequest.Permissions)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	role := role.Role{
		ID:          uuid.New(),
		Name:        validatedRequest.Name,
		Permissions: permissions,
	}
	if err = useCase.roleRepository.Save(role); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (useCase *CreateRoleUseCase) findPermissions(permissionNames []string) ([]permission.Permission, error) {
	permissions, err := useCase.permissionRepository.FindByNames(permissionNames)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (*CreateRoleUseCase) RequiredPermissions() []string {
	return []string{role.CreateRolePermission}
}

func NewCreateRoleUseCase(roleRepository role.RoleRepository, permissionRepository permission.PermissionRepository, logger *zap.Logger) *CreateRoleUseCase {
	useCase := CreateRoleUseCase{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		logger:               logger,
	}
	return &useCase
}
