// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"crypto/tls"
	"github.com/elkopass/BITA/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func createClientConn() (*grpc.ClientConn, error) {
	tlsConfig := tls.Config{}

	return grpc.Dial(config.ApiURL, grpc.WithTransportCredentials(credentials.NewTLS(&tlsConfig)))
}
