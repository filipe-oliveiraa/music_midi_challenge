package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"crossjoin.com/gorxestra/config"
	musiciand "crossjoin.com/gorxestra/daemon/musiciand"
	"crossjoin.com/gorxestra/logging"
	"crossjoin.com/gorxestra/util/conf"
	"github.com/gofrs/flock"
)

const (
	Success = 0
	Failed  = 1
)

const (
	LockFile = "conductor.lock"
)

func main() {
	var cfg config.MusicianConf
	help, err := conf.ParseConfig(&cfg)
	if errors.Is(err, conf.ErrHelp) {
		fmt.Println(help)
		return
	}

	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(Failed)
	}
}

func run(cfg config.MusicianConf) error {
	absolutePath, absPathErr := filepath.Abs(cfg.DataDir)
	if len(cfg.DataDir) == 0 {
		return fmt.Errorf("data directory not specified. Please use -d or set in your environment")
	}

	if absPathErr != nil {
		return fmt.Errorf("can't convert data directory's path to absolute, %v", cfg.DataDir)
	}

	if _, err := os.Stat(absolutePath); err != nil {
		return err
	}

	// before doing anything further, attempt to acquire the service lock
	// to ensure this is the only node running against this data directory
	lockPath := filepath.Join(absolutePath, LockFile)
	fileLock := flock.New(lockPath)
	locked, err := fileLock.TryLock()

	log := logging.Base()

	if err != nil {
		return fmt.Errorf("unexpected failure in establishing %s: %w", LockFile, err)
	}

	if !locked {
		return fmt.Errorf("failed to lock %s; an instance of the service already running in this data directory", LockFile)
	}
	defer fileLock.Unlock() //nolint: errcheck

	s := musiciand.Server{
		RootPath: absolutePath,
	}

	if cfg.Logger.LogToStdout {
		cfg.Logger.LogSizeLimit = 0
	}

	err = s.Initialize(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.Errorf("initializing service: %w", err)
		return fmt.Errorf("initializing service: %w", err)
	}

	if cfg.InitAndExit {
		return nil
	}

	s.Start()
	return nil
}
