package credentials

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func GetServerOptions(certChain *CertChain) ([]grpc.ServerOption, error) {
	opts := []grpc.ServerOption{}
	if certChain == nil {
		return nil, nil
	}

	cp := x509.NewCertPool()
	cp.AppendCertsFromPEM(certChain.RootCA)

	cert, err := tls.X509KeyPair(certChain.Cert, certChain.Key)
	if err != nil {
		return nil, err
	}

	//nolint:gosec
	config := &tls.Config{
		ClientCAs: cp,
		// Require cert verification
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	opts = append(opts, grpc.Creds(credentials.NewTLS(config)))

	return opts, nil
}

func GetClientOptions(certChain *CertChain, serverName string) ([]grpc.DialOption, error) {
	opts := []grpc.DialOption{}
	if certChain != nil {
		cp := x509.NewCertPool()
		ok := cp.AppendCertsFromPEM(certChain.RootCA)
		if !ok {
			return nil, errors.New("failed to append PEM root cert to x509 CertPool")
		}
		config, err := TLSConfigFromCertAndKey(certChain.Cert, certChain.Key, serverName, cp)
		if err != nil {
			return nil, fmt.Errorf("failed to create tls config from cert and key: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return opts, nil
}
