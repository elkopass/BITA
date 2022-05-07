package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	pb "github.com/elkopass/BITA/internal/proto"
)

type OrdersStreamInterface interface {
	// Stream сделок пользователя
	TradesStream(in *pb.TradesStreamRequest) (pb.OrdersStreamService_TradesStreamClient, error)
}

type OrdersStreamService struct {
	client pb.OrdersStreamServiceClient
}

func NewOrdersStreamService() *OrdersStreamService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewOrdersStreamServiceClient(conn)
	return &OrdersStreamService{client: client}
}

func (oss OrdersStreamService) TradesStream(in *pb.TradesStreamRequest) (pb.OrdersStreamService_TradesStreamClient, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := oss.client.TradesStream(ctx, in)
	if err != nil {
		return nil, err
	}

	return res, nil
}
