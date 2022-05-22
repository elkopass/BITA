// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"context"
	"fmt"
	"github.com/elkopass/BITA/internal/config"
	"github.com/google/uuid"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// createRequestContext returns context for API calls with timeout and auth headers attached.
func createRequestContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultRequestTimeout)

	authHeader := fmt.Sprintf("Bearer %s", config.TradeBotConfig().Token)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "x-tracking-id", uuid.New().String())
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "x-app-name", config.AppName)

	return ctx, cancel
}

// createRequestContext returns context for streams with auth headers attached.
func createStreamContext() context.Context {
	ctx := context.TODO()

	authHeader := fmt.Sprintf("Bearer %s", config.TradeBotConfig().Token)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "x-tracking-id", uuid.New().String())
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "x-app-name", config.AppName)

	return ctx
}
