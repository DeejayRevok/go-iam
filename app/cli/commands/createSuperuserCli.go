package commands

import (
	"context"
	"go-iam/src/application/createUser"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type CreateSuperuserCLI struct {
	createUserUseCase *createUser.CreateUserUseCase
	logger            *zap.Logger
}

func (cli *CreateSuperuserCLI) Execute(c *cli.Context) error {
	cli.logger.Info("Starting superuser creation")
	defer cli.logger.Info("Finished superuser creation")

	username := c.String("username")
	email := c.String("email")
	password := c.String("password")

	ctx := context.Background()

	superuserRequest := createUser.CreateUserRequest{
		Username:  username,
		Email:     email,
		Password:  password,
		Superuser: true,
	}
	useCaseResponse := cli.createUserUseCase.Execute(ctx, &superuserRequest)
	return useCaseResponse.Err
}

func NewCreateSuperuserCLI(createUserUseCase *createUser.CreateUserUseCase, logger *zap.Logger) *CreateSuperuserCLI {
	return &CreateSuperuserCLI{
		createUserUseCase: createUserUseCase,
		logger:            logger,
	}
}
