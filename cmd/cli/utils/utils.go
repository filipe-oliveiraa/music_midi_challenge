package utils

import (
	conductor "crossjoin.com/gorxestra/daemon/conductord/api/client/v1"
	"github.com/urfave/cli/v2"
)

const ConductorAddressFlag = "conductorAddr"

func GetConductorCli(ctx *cli.Context) (conductor.ClientDaemon, error) {
	addr := ctx.String(ConductorAddressFlag)
	cli, err := conductor.New(addr)
	if err != nil {
		return nil, err
	}
	return cli, err
}

/*
func getMusicianCli(ctx cli.Context) (musician.ClientDaemon, error) {
	cli, err := musician.New()
	if err != nil {
		return nil, err
	}
	return cli, err
}*/
