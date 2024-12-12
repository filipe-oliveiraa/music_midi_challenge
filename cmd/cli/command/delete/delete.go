package delete

import (
	"errors"

	"crossjoin.com/gorxestra/cmd/cli/utils"
	"crossjoin.com/gorxestra/data"
	"github.com/urfave/cli/v2"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:         "delete",
		Aliases:      nil,
		Usage:        "<musician ID>",
		UsageText:    "",
		Description:  "Delete a musician",
		Args:         false,
		ArgsUsage:    "",
		Category:     "Basic Commands (Beginner)",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action:       deleteAction,
		OnUsageError: nil,
		Subcommands:  cli.Commands{},
		//nolint
		Flags:                  []cli.Flag{},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}
}

func deleteAction(ctx *cli.Context) error {
	idRaw := ctx.Args().First()
	if idRaw == "" {
		return errors.New("specify a id")
	}

	id, err := data.IdFromHex(idRaw)
	if err != nil {
		return errors.New("invalid musician id")
	}

	cli, err := utils.GetConductorCli(ctx)
	if err != nil {
		return err
	}

	return cli.UnregisterMusician(id)
}
