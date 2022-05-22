// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	pb "github.com/elkopass/BITA/internal/proto"
)

type OrdersStreamInterface interface {
	// Recv listens for incoming messages and block until first one is received.
	Recv() (*pb.TradesStreamResponse, error)
}

type OrdersStream struct {
	client pb.OrdersStreamServiceClient
	stream pb.OrdersStreamService_TradesStreamClient
}

func NewOrdersStream(request *pb.TradesStreamRequest) *OrdersStream{
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewOrdersStreamServiceClient(conn)
	ctx := createStreamContext()

	stream, err := client.TradesStream(ctx, request)
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	return &OrdersStream{client: client, stream: stream}
}

func (os OrdersStream) Recv() (*pb.TradesStreamResponse, error) {
	return os.stream.Recv()
}
