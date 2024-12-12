package broker

import (
	"context"
	"os"
	"path"

	"crossjoin.com/gorxestra/config"
	"crossjoin.com/gorxestra/data"
	"crossjoin.com/gorxestra/logging"
	"crossjoin.com/gorxestra/service/conductor/baton"
)

type ConductorNode struct {
	log     logging.Logger
	rootDir string

	config config.ConductorConf

	baton baton.Baton

	ctx    context.Context
	cancel context.CancelFunc
}

func New(log logging.Logger, rootDir string, cfg config.ConductorConf) (*ConductorNode, error) {
	ctx, cancel := context.WithCancel(context.Background())
	c := ConductorNode{
		log:     log,
		rootDir: rootDir,
		config:  cfg,
		baton:   baton.New(log),
		ctx:     ctx,
		cancel:  cancel,
	}

	return &c, nil
}

func (c *ConductorNode) RegisterMusician(m data.Musician) error {
	c.log.
		With("id", m.Id.Hex()).
		Info("registering musician")
	return c.baton.RegisterMusician(m)
}

func (c *ConductorNode) UnregisterMusician(id data.ID) error {
	c.log.
		With("id", id.Hex()).
		Info("registering musician")
	return c.baton.UnregisterMusician(id)
}

func (c *ConductorNode) PlayMusic(name string) error {
	path := path.Join(c.rootDir, name)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	c.baton.Play(f)

	return nil

}

func (c *ConductorNode) Start() error {
	return nil
}

func (broker *ConductorNode) Stop() error {
	return nil
}

func (broker *ConductorNode) Config() config.ConductorConf {
	return broker.config
}

func (broker *ConductorNode) Status() error {
	return nil
}
