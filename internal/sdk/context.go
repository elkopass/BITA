package sdk

import (
	"context"
	"fmt"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/config"
	"github.com/google/uuid"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func createRequestContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)

	authHeader := fmt.Sprintf("Bearer %s", config.TradeBotConfig().Token)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "x-tracking-id", uuid.New().String())
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "x-app-name", config.AppName)

	return ctx, cancel
}
