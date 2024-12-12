package connection

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"
)

type tlsCtxKey uint8

const (
	PeerCertificates tlsCtxKey = iota
)

type tlsListener struct {
	listener net.Listener
}

func newTlsListener(listener net.Listener) Listener {
	return &tlsListener{
		listener: listener,
	}
}

// Accept implements Listener.
func (t *tlsListener) Accept() (Conn, error) {
	conn, err := t.listener.Accept()
	if err != nil {
		return nil, err
	}

	return newTlsConn(conn)
}

// Addr implements Listener.
func (t *tlsListener) Addr() net.Addr {
	return t.listener.Addr()
}

// Close implements Listener.
func (t *tlsListener) Close() error {
	return t.listener.Close()
}

type tlsConn struct {
	*tls.Conn
	once sync.Once
	ctx  context.Context
}

func newTlsConn(conn net.Conn) (Conn, error) {
	tls, ok := conn.(*tls.Conn)

	if !ok {
		return nil, errors.New("not a tls connection")
	}

	return &tlsConn{
		Conn: tls,
		once: sync.Once{},
		ctx:  context.Background(),
	}, nil
}

func (c *tlsConn) Context() context.Context {
	c.once.Do(c.extractContext)
	return c.ctx
}

func (c *tlsConn) extractContext() {
	ctx := c.ctx
	_ = c.Conn.Handshake()
	state := c.Conn.ConnectionState()

	ctx = context.WithValue(ctx, PeerCertificates, state.PeerCertificates)
	c.ctx = ctx
}

type identityListener struct {
	listener net.Listener
}

func newIdentityListener(listener net.Listener) Listener {
	return &identityListener{
		listener: listener,
	}
}

// Accept implements Listener.
func (t *identityListener) Accept() (Conn, error) {
	conn, err := t.listener.Accept()
	if err != nil {
		return nil, err
	}
	return newIdentityConn(conn)
}

// Addr implements Listener.
func (t *identityListener) Addr() net.Addr {
	return t.listener.Addr()
}

// Close implements Listener.
func (t *identityListener) Close() error {
	return t.listener.Close()
}

type identityConn struct {
	net.Conn
}

func newIdentityConn(conn net.Conn) (Conn, error) {
	return &identityConn{
		Conn: conn,
	}, nil
}

func (c *identityConn) Context() context.Context {
	return context.Background()
}
