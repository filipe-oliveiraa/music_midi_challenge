package add

import (
	"errors"

	"crossjoin.com/gorxestra/cmd/cli/utils"
	"crossjoin.com/gorxestra/data"
	"github.com/urfave/cli/v2"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:         "add",
		Aliases:      nil,
		Usage:        "<musician ID> <musician addr>",
		UsageText:    "",
		Description:  "Add a musician",
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
	if ctx.Args().Len() != 2 {
		return errors.New("please specify: <musician ID> <musician addr>")
	}

	idRaw := ctx.Args().Get(0)
	if idRaw == "" {
		return errors.New("specify a id")
	}

	id, err := data.IdFromHex(idRaw)
	if err != nil {
		return errors.New("invalid musician id")
	}

	addr := ctx.Args().Get(1)
	if addr == "" {
		return errors.New("specify a musician url")
	}

	cli, err := utils.GetConductorCli(ctx)
	if err != nil {
		return err
	}

	return cli.RegisterMusician(data.Musician{
		Id:      id,
		Address: addr,
	})
}
