package commands

import (
	"go-uaa/src/application/createPermission"
	"go-uaa/src/domain/permission"
	"go-uaa/src/domain/role"
	"go-uaa/src/domain/user"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var permissions [8]string = [...]string{
	permission.CreatePermissionPermission,
	role.CreateRolePermission,
	role.UpdateRolePermission,
	role.DeleteRolePermission,
	user.CreateUserPermission,
	user.DeleteUserPermission,
	user.ReadUserPermission,
	user.UpdateUserPermission,
}

type BoostrapPermissionsCLI struct {
	createPermissionUseCase *createPermission.CreatePermissionUseCase
	logger                  *zap.Logger
}

func (cli *BoostrapPermissionsCLI) Execute(_ *cli.Context) error {
	cli.logger.Info("Starting permissions bootstraping")
	defer cli.logger.Info("Finished permissions bootstraping")
	for _, permission := range permissions {
		useCaseRequest := createPermission.CreatePermissionRequest{
			Name: permission,
		}
		response := cli.createPermissionUseCase.Execute(&useCaseRequest)
		if response.Err != nil {
			return response.Err
		}
	}
	return nil
}

func NewBoostrapPermissionsCLI(createPermissionUseCase *createPermission.CreatePermissionUseCase, logger *zap.Logger) *BoostrapPermissionsCLI {
	return &BoostrapPermissionsCLI{
		createPermissionUseCase: createPermissionUseCase,
		logger:                  logger,
	}
}
