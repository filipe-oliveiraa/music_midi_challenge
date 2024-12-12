package main

import (
	"context"
	"fmt"
	"os"

	"crossjoin.com/gorxestra/cmd/cli/command"
	"crossjoin.com/gorxestra/cmd/cli/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	//nolint
	app := &cli.App{
		Name:  "pipernet",
		Usage: "pipernet client",
		Flags: []cli.Flag{
			//nolint
			&cli.StringFlag{
				Name:    utils.ConductorAddressFlag,
				Value:   "http://localhost:8080",
				Aliases: []string{"c"},
			},
		},
		Commands: command.GetCommands(),
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		return err
	}

	return nil
}
