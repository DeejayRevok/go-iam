package main

import (
	"go-uaa/app"
	"go-uaa/app/cli/commands"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var clis []*cli.Command

	container := app.BuildDIContainer()
	container.Invoke(func(permissionsCli *commands.BoostrapPermissionsCLI) {
		clis = append(clis, &cli.Command{
			Name:   "BootstrapPermissions",
			Usage:  "Bootstrap the application permissions",
			Action: permissionsCli.Execute,
		})
	})
	container.Invoke(func(createSuperuserCli *commands.CreateSuperuserCLI) {
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
	})

	app := &cli.App{
		Commands: clis,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
