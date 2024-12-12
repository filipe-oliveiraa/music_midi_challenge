package hellod

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"crossjoin.com/gorxestra/config"
	apiServer "crossjoin.com/gorxestra/daemon/conductord/api/server"
	"crossjoin.com/gorxestra/logging"
	broker "crossjoin.com/gorxestra/service/conductor"
	"github.com/labstack/echo/v4"

	"crossjoin.com/gorxestra/util/network/limitlistener"

	"crossjoin.com/gorxestra/util"
)

var server http.Server

// maxHeaderBytes must have enough room to hold an api token
const (
	maxHeaderBytes = 4096
	pidFileName    = "conductor.pid"
	httpFileName   = "conductor.http"
)

// ServerNode is the required methods for any node the server fronts
type ServerNode interface {
	apiServer.APINodeInterface
	Start() error
	Stop() error
}

// Server represents an instance of the REST API HTTP server
type Server struct {
	RootPath string

	pidFile  string
	httpFile string

	httpListener net.Listener

	log      logging.Logger
	node     ServerNode
	stopping chan struct{}
}

// Initialize creates a Node instance with applicable network services
func (s *Server) Initialize(cfg config.ConductorConf) error {
	// set up node
	s.log = logging.Base()

	s.log.SetOutput(os.Stdout)
	s.log.SetLevel(logging.Level(cfg.Logger.BaseLoggerDebugLevel))

	node, err := broker.New(s.log, s.RootPath, cfg)

	if os.IsNotExist(err) {
		return fmt.Errorf("node has not been installed: %s", err)
	}
	if err != nil {
		return fmt.Errorf("couldn't initialize the node: %s", err)
	}

	s.node = ServerNode(node)

	// When a caller to logging uses Fatal, we want to stop the node before os.Exit is called.
	logging.RegisterExitHandler(s.Stop)

	return nil
}

// Start starts a Node instance and its network services
func (s *Server) Start() {
	s.log.Info("Trying to start a Conductor")
	_ = s.node.Start()
	s.log.Info("Successfully started a Conductor.")

	cfg := s.node.Config()

	s.stopping = make(chan struct{})

	if err := s.setPidFile(); err != nil {
		os.Exit(1)
	}

	// Handle signals cleanly
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGHUP)

	errChan := make(chan error, 1)

	s.newHttpServer(cfg, errChan)

	s.log.With("httpAddr", s.httpListener.Addr().String()).
		Info("Conductor running and accepting requests over HTTP. Press Ctrl-C to exit.")

	select {
	case err := <-errChan:
		if err != nil {
			s.log.Warn(err)
		} else {
			s.log.Info("Node exited successfully")
		}
		s.Stop()
	case sig := <-c:
		fmt.Printf("Exiting on %v\n", sig)

		s.Stop()
		os.Exit(0)
	}
}

func (s *Server) newHttpServer(cfg config.ConductorConf, errChan chan error) (*echo.Echo, string) {
	listener, err := util.MakeListener(cfg.Rest.EndpointAddress)
	if err != nil {
		fmt.Printf("Could not start node: %v\n", err)
		os.Exit(1)
	}

	listener = limitlistener.RejectingLimitListener(
		listener, cfg.Rest.ConnectionsHardLimit, s.log)

	addr := listener.Addr().String()
	//nolint: exhaustruct
	server = http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(cfg.Rest.ReadTimeoutSeconds) * time.Second,
		WriteTimeout:   time.Duration(cfg.Rest.WriteTimeoutSeconds) * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}

	s.httpListener = listener

	e := apiServer.NewHttpRouter(
		s.log,
		s.node,
		s.stopping,
		listener,
		cfg.Rest.ConnectionsSoftLimit,
	)

	go func() {
		err := e.StartServer(&server)
		errChan <- err
	}()

	if err := s.setHttpFile(addr); err != nil {
		os.Exit(1)
	}

	return e, addr
}

func (s *Server) setPidFile() error {
	s.pidFile = filepath.Join(s.RootPath, pidFileName)
	err := os.WriteFile(s.pidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0o600)
	if err != nil {
		fmt.Printf("pidfile error: %v\n", err)
	}

	return err
}

func (s *Server) setHttpFile(addr string) error {
	s.httpFile = filepath.Join(s.RootPath, httpFileName)
	err := os.WriteFile(s.httpFile, []byte(fmt.Sprintf("%s\n", addr)), 0o600)
	if err != nil {
		fmt.Printf("netfile error: %v\n", err)
		return err
	}
	return err
}

// Stop initiates a graceful shutdown of the node by shutting down the network server.
func (s *Server) Stop() {
	// close the s.stopping, which would signal the rest api router that any pending commands
	// should be aborted.
	close(s.stopping)

	s.httpListener.Close()

	err := s.node.Stop()
	if err != nil {
		s.log.Error(err)
	}

	err = server.Shutdown(context.Background())
	if err != nil {
		s.log.Error(err)
	}

	os.Remove(s.pidFile)
	os.Remove(s.httpFile)
}
