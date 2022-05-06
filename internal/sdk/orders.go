package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/metrics"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type OrdersInterface interface {
	// Метод выставления заявки.
	PostOrder(order *pb.PostOrderRequest) (*pb.PostOrderResponse, error)
	// Метод отмены биржевой заявки.
	CancelOrder(accountID string, orderID string) (*timestamp.Timestamp, error)
	// Метод получения статуса торгового поручения.
	GetOrderState(accountID string, orderID string) (*pb.OrderState, error)
	// Метод получения списка активных заявок по счёту.
	GetOrders(accountID string) ([]*pb.OrderState, error)
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

	os.incrementRequestsCounter("PostOrder")
	res, err := os.client.PostOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (os OrdersService) CancelOrder(accountID string, orderID string) (*timestamp.Timestamp, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	os.incrementRequestsCounter("CancelOrder")
	res, err := os.client.CancelOrder(ctx, &pb.CancelOrderRequest{
		AccountId: accountID,
		OrderId:   orderID,
	})
	if err != nil {
		return nil, err
	}

	return res.Time, nil
}

func (os OrdersService) GetOrderState(accountID string, orderID string) (*pb.OrderState, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	os.incrementRequestsCounter("GetOrderState")
	res, err := os.client.GetOrderState(ctx, &pb.GetOrderStateRequest{
		AccountId: accountID,
		OrderId:   orderID,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (os OrdersService) GetOrders(accountID string) ([]*pb.OrderState, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	os.incrementRequestsCounter("GetOrders")
	res, err := os.client.GetOrders(ctx, &pb.GetOrdersRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, err
	}

	return res.Orders, nil
}

func (os OrdersService) incrementRequestsCounter(method string) {
	metrics.ApiRequests.WithLabelValues(loggy.GetBotID(), "OrdersService", method).Inc()
}
