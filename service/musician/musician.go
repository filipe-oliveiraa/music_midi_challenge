package musician

import (
	"context"
	"fmt"
	"time"

	"crossjoin.com/gorxestra/config"
	"crossjoin.com/gorxestra/daemon/conductord/api/client/v1"
	"crossjoin.com/gorxestra/data"
	"crossjoin.com/gorxestra/logging"
	"gitlab.com/gomidi/midi/v2"

	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv" // autoregisters driver
)

type MusicianNode struct {
	log     logging.Logger
	rootDir string

	config config.MusicianConf
	cli    client.ClientDaemon

	out    drivers.Out
	ctx    context.Context
	cancel context.CancelFunc
}

func New(log logging.Logger, rootDir string, cfg config.MusicianConf) (*MusicianNode, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cli, err := client.New(cfg.Conductor.ConductorAddr)

	if err != nil {
		return nil, err
	}

	m := MusicianNode{
		out:     nil,
		log:     log,
		rootDir: rootDir,
		config:  cfg,
		cli:     cli,
		ctx:     ctx,
		cancel:  cancel,
	}

	return &m, nil
}

func (m *MusicianNode) Play(bs []byte) error {
	m.log.With("note", bs).Info("playing sound")
	return m.out.Send(bs)
}

func (m *MusicianNode) registerMusician() error {

	for {
		m.log.Info("attemp to register node")

		err := m.cli.RegisterMusician(data.Musician{
			Id:      data.GenId(),
			Address: m.config.Conductor.AdvertiseAddr,
		})
		if err == nil {
			return nil
		}

		m.log.With("error", err).Warn("attempting again to connect")
		time.Sleep(1 * time.Second)
	}
}

func (m *MusicianNode) unregisterMusician() error {
	//TODO currently not implemented. Not needed for the challenge
	return nil
}

func (m *MusicianNode) Start() error {
	out, err := midi.OutPort(1) //midi.FindOutPort("qsynth")
	if err != nil {
		return fmt.Errorf("could not open midi port:%w", err)
	}

	m.out = out

	return m.registerMusician()
}

func (m *MusicianNode) Stop() error {
	midi.CloseDriver()
	return nil
}

func (m *MusicianNode) Config() config.MusicianConf {
	return m.config
}

func (m *MusicianNode) Status() error {
	return nil
}
