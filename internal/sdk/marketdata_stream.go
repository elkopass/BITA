// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	pb "github.com/elkopass/BITA/internal/proto"
)

type MarketDataStreamInterface interface {
	// Bi-directional стрим предоставления биржевой информации.
	MarketDataStream() (pb.MarketDataStreamService_MarketDataStreamClient, error)
}

type MarketDataStreamService struct {
	client pb.MarketDataStreamServiceClient
}

func NewMarketDataStreamService() *MarketDataStreamService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewMarketDataStreamServiceClient(conn)
	return &MarketDataStreamService{client: client}
}

func (mdss MarketDataStreamService) MarketDataStream() (pb.MarketDataStreamService_MarketDataStreamClient, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := mdss.client.MarketDataStream(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
