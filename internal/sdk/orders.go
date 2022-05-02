package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type OrdersStreamInterface interface {
	// Stream сделок пользователя
	TradesStream(in *pb.TradesStreamRequest) (pb.OrdersStreamService_TradesStreamClient, error)
}

type OrdersInterface interface {
	// Метод выставления заявки.
	PostOrder(order *pb.PostOrderRequest) (*pb.PostOrderResponse, error)
	// Метод отмены биржевой заявки.
	CancelOrder(accountID AccountID, orderID OrderID) (*timestamp.Timestamp, error)
	// Метод получения статуса торгового поручения.
	GetOrderState(accountID AccountID, orderID OrderID) (*pb.OrderState, error)
	// Метод получения списка активных заявок по счёту.
	GetOrders(accountID AccountID) ([]*pb.OrderState, error)
}

type OrdersService struct {
	client pb.OrdersServiceClient
}

func NewOrdersService() *OrdersService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewOrdersServiceClient(conn)
	return &OrdersService{client: client}
}

func (os OrdersService) PostOrder(order *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := os.client.PostOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (os OrdersService) CancelOrder(accountID AccountID, orderID OrderID) (*timestamp.Timestamp, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := os.client.CancelOrder(ctx, &pb.CancelOrderRequest{
		AccountId: string(accountID),
		OrderId: string(orderID),
	})
	if err != nil {
		return nil, err
	}

	return res.Time, nil
}

func (os OrdersService) GetOrderState(accountID AccountID, orderID OrderID) (*pb.OrderState, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := os.client.GetOrderState(ctx, &pb.GetOrderStateRequest{
		AccountId: string(accountID),
		OrderId: string(orderID),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (os OrdersService) GetOrders(accountID AccountID) ([]*pb.OrderState, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := os.client.GetOrders(ctx, &pb.GetOrdersRequest{
		AccountId: string(accountID),
	})
	if err != nil {
		return nil, err
	}

	return res.Orders, nil
}
