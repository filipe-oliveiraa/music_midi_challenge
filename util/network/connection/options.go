package connection

import (
	"crypto/tls"
	"crypto/x509"
)

type TLSOption func(*tlsOptions) error

type tlsOptions struct {
	certs              []tls.Certificate
	caClient           *x509.CertPool
	rootCA             *x509.CertPool
	clientAuthType     tls.ClientAuthType
	insecureSkipVerify bool
	serverName         string
}

func newListenOptions() tlsOptions {
	return tlsOptions{
		certs:              nil,
		caClient:           x509.NewCertPool(),
		rootCA:             x509.NewCertPool(),
		clientAuthType:     tls.NoClientCert,
		serverName:         "localhost",
		insecureSkipVerify: false,
	}
}

func getListenOptions(opts []TLSOption) (tlsOptions, error) {
	o := newListenOptions()
	for i := range opts {
		if err := opts[i](&o); err != nil {
			return o, err
		}
	}

	return o, nil
}

func WithRootCA(cert *x509.Certificate) TLSOption {
	return func(l *tlsOptions) error {
		l.caClient.AddCert(cert)
		return nil
	}
}

// WithCertificate includes
func WithCertificate(cert, key []byte) TLSOption {
	return func(l *tlsOptions) error {
		cert, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return err
		}
		l.certs = append(l.certs, cert)
		return nil
	}
}

func WithInsecure(insecure bool) TLSOption {
	return func(l *tlsOptions) error {
		l.insecureSkipVerify = insecure
		return nil
	}
}

func WithClientCA(cert *x509.Certificate) TLSOption {
	return func(l *tlsOptions) error {
		l.caClient.AddCert(cert)
		return nil
	}
}

func WithEnabledMTLS() TLSOption {
	return WithMTLS(true)
}

func WithMTLS(enable bool) TLSOption {
	return func(l *tlsOptions) error {
		if enable {
			l.clientAuthType = tls.RequireAndVerifyClientCert
		} else {
			l.clientAuthType = tls.NoClientCert
		}
		return nil
	}
}
