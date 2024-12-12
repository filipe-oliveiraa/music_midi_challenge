package command

import (
	"crossjoin.com/gorxestra/cmd/cli/command/add"
	"crossjoin.com/gorxestra/cmd/cli/command/delete"
	"crossjoin.com/gorxestra/cmd/cli/command/play"
	"github.com/urfave/cli/v2"
)

func GetCommands() cli.Commands {
	return cli.Commands{
		play.Commands(),
		delete.Commands(),
		add.Commands(),
	}
}
