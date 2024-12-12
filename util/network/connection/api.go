package connection

import (
	"context"
	"net"
)

// Extended Listener
type Listener interface {
	// Accept waits for and returns the next connection to the listener
	// and a context
	Accept() (Conn, error)

	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error

	// Addr returns the listener's network address.
	Addr() net.Addr
}

// Extended Connection
type Conn interface {
	net.Conn
	Context() context.Context
}
