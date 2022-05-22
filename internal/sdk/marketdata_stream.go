// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	pb "github.com/elkopass/BITA/internal/proto"
)

type MarketDataStreamInterface interface {
	// Recv listens for incoming messages and block until first one is received.
	Recv() (*pb.MarketDataResponse, error)
	// Send puts pb.MarketDataRequest into a stream.
	Send(request *pb.MarketDataRequest) error
}

type MarketDataStream struct {
	client pb.MarketDataStreamServiceClient
	stream pb.MarketDataStreamService_MarketDataStreamClient
}

func NewMarketDataStream() *MarketDataStream {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewMarketDataStreamServiceClient(conn)
	ctx := createStreamContext()

	stream, err := client.MarketDataStream(ctx)
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	return &MarketDataStream{client: client, stream: stream}
}

func (mds MarketDataStream) Recv() (*pb.MarketDataResponse, error) {
	return mds.stream.Recv()
}

func (mds MarketDataStream) Send(request *pb.MarketDataRequest) error {
	return mds.stream.Send(request)
}
