package main

import (
	"go-iam/app"
	"go-iam/app/cli/commands"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	var clis []*cli.Command

	container := app.BuildDIContainer()
	if err := container.Invoke(func(logger *zap.Logger) {
		handleError(container.Invoke(func(createSuperuserCli *commands.CreateSuperuserCLI) {
			clis = append(clis, &cli.Command{
				Name:   "CreateSuperuser",
				Usage:  "Create a new superuser",
				Action: createSuperuserCli.Execute,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "username",
						Aliases:  []string{"u"},
						Usage:    "Username of the superuser",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "email",
						Aliases:  []string{"e"},
						Usage:    "Email of the superuser",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Aliases:  []string{"p"},
						Usage:    "Password of the superuser",
						Required: true,
					},
				},
			})
		}), logger)
	}); err != nil {
		panic("Error trying to build command clis")
	}

	app := &cli.App{
		Commands: clis,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func handleError(err error, logger *zap.Logger) {
	if err != nil {
		logger.Fatal(err.Error())
	}
}
