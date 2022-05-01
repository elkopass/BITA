package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/config"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
	"time"
)


func main() {
	log := loggy.GetLogger().Sugar()

	tradeBotConfig := config.GetTradeBotConfig()
	if tradeBotConfig.Token == "" {
		log.Fatal("TRADEBOT_TOKEN environment variable is required to run this program")
	}

	tlsConfig := tls.Config{}
	conn, err := grpc.Dial(config.ApiURL, grpc.WithTransportCredentials(credentials.NewTLS(&tlsConfig)))
	if err != nil {
		log.Fatalf("did not connect: %v", err.Error())
	} else {
		defer conn.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	authHeader := fmt.Sprintf("Bearer %s", tradeBotConfig.Token)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "x-tracking-id", uuid.New().String())
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "x-app-name", config.AppName)

	sandboxServiceClient := pb.NewSandboxServiceClient(conn)
	r1, err := sandboxServiceClient.GetSandboxAccounts(ctx, &pb.GetAccountsRequest{})
	if err != nil {
		log.Error(err.Error())
	}

	log.Info(r1)
}
