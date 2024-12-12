package play

import (
	"errors"

	"crossjoin.com/gorxestra/cmd/cli/utils"
	"github.com/urfave/cli/v2"
)

func Commands() *cli.Command {
	return &cli.Command{
		Name:         "play",
		Aliases:      nil,
		Usage:        "<musician>",
		UsageText:    "",
		Description:  "Play a music",
		Args:         false,
		ArgsUsage:    "",
		Category:     "Basic Commands (Beginner)",
		BashComplete: nil,
		Before:       nil,
		After:        nil,
		Action:       playAction,
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

func playAction(ctx *cli.Context) error {
	music := ctx.Args().First()
	if music == "" {
		return errors.New("specify a music")
	}

	cli, err := utils.GetConductorCli(ctx)
	if err != nil {
		return err
	}

	return cli.PlayMusic(music)
}
