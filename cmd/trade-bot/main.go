package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/config"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
	"time"
)


func main() {
	log := loggy.GetLogger().Sugar()

	var tradeBotConfig config.TradeBotConfig

	err := envconfig.Process("tradebot", &tradeBotConfig)
	if err != nil {
		log.Fatalf("failed to process config: %v", err)
	}

	if tradeBotConfig.Token == "" {
		log.Fatal("TRADEBOT_TOKEN environment variable is required to run this program")
	}

	tlsConfig := tls.Config{}
	conn, err := grpc.Dial(tradeBotConfig.ApiURL, grpc.WithTransportCredentials(credentials.NewTLS(&tlsConfig)))
	if err != nil {
		log.Fatalf("did not connect: %v", err.Error())
	} else {
		defer conn.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	authHeader := fmt.Sprintf("Bearer %s", tradeBotConfig.Token)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", authHeader)

	client := pb.NewSandboxServiceClient(conn)
	r, err := client.GetSandboxAccounts(ctx, &pb.GetAccountsRequest{})
	if err != nil {
		log.Error(err.Error())
	}

	log.Info(r)
}
