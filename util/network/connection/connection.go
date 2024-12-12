package connection

import (
	"crypto/tls"
	"net"
	"strings"
	"time"
)

const (
	UTP = "utp"
	TCP = "tcp"
)

type DialerConf struct {
	Protocol        string
	Address         string
	Timeout         time.Duration
	ReadBufferSize  int
	WriteBufferSize int
}

// Dial wraps around the go net.Dial and adds utp protocol
func Dial(dialer DialerConf) (Conn, error) {
	conn, err := dial(dialer)
	if err != nil {
		return nil, err
	}

	return newIdentityConn(conn)
}

func dial(dialer DialerConf) (net.Conn, error) {
	var conn net.Conn
	var err error

	switch dialer.Protocol {
	default:
		// nolint
		dial := net.Dialer{Timeout: dialer.Timeout}
		conn, err = dial.Dial(dialer.Protocol, dialer.Address)
	}

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// DialTLS adds TLS to the Dial
func DialTLS(dialer DialerConf, opts ...TLSOption) (Conn, error) {
	conn, err := dial(dialer)
	if err != nil {
		return nil, err
	}

	o, err := getListenOptions(opts)
	if err != nil {
		return nil, err
	}

	//nolint:exhaustruct
	connTls := tls.Client(conn, &tls.Config{
		RootCAs:      o.rootCA,
		Certificates: o.certs,
		// nolint:gosec
		InsecureSkipVerify: o.insecureSkipVerify,
		ServerName:         o.serverName,
	})

	return newTlsConn(connTls)
}

func Listen(dialer DialerConf) (Listener, error) {
	listener, err := listenWithConversion(dialer)
	if err != nil {
		return nil, err
	}
	return newIdentityListener(listener), err
}

func listenWithConversion(dialer DialerConf) (net.Listener, error) {
	var listener net.Listener
	var err error
	if (dialer.Address == "127.0.0.1:0") || (dialer.Address == ":0") {
		// if port 0 is provided, prefer port 8080 first, then fall back to port 0
		preferredAddr := strings.Replace(dialer.Address, ":0", ":8080", -1)
		listener, err = listen(dialer.Protocol, preferredAddr)
		if err == nil {
			return listener, err
		}
	}
	// err was not nil or :0 was not provided, fall back to originally passed addr
	return listen(dialer.Protocol, dialer.Address)
}

func listen(network, addr string) (net.Listener, error) {
	var listener net.Listener
	var err error
	switch network {
	default:
		listener, err = net.Listen(network, addr)
	}

	if err != nil {
		return nil, err
	}

	return listener, nil
}

func ListenTLS(dialer DialerConf, opts ...TLSOption) (Listener, error) {
	l, err := listenWithConversion(dialer)
	if err != nil {
		return nil, err
	}

	o, err := getListenOptions(opts)
	if err != nil {
		return nil, err
	}

	// nolint:exhaustruct
	tl := tls.NewListener(l, &tls.Config{
		// Server Options
		ClientCAs:    o.caClient,
		ClientAuth:   o.clientAuthType,
		Certificates: o.certs,
		// nolint:gosec
		InsecureSkipVerify: o.insecureSkipVerify,
	})

	return newTlsListener(tl), nil
}
