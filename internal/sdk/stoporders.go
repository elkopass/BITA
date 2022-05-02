package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type StopOrdersInterface interface {
	// Метод выставления стоп-заявки.
	PostStopOrder(stopOrder *pb.PostStopOrderRequest) (StopOrderID, error)
	// Метод получения списка активных стоп заявок по счёту.
	GetStopOrders(accountID AccountID) ([]*pb.StopOrder, error)
	// Метод отмены стоп-заявки.
	CancelStopOrder(accountID AccountID, stopOrderID StopOrderID) (*timestamp.Timestamp, error)
}

type StopOrdersService struct {
	client pb.StopOrdersServiceClient
}

func NewStopOrdersService() *StopOrdersService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewStopOrdersServiceClient(conn)
	return &StopOrdersService{client: client}
}

func (sos StopOrdersService) PostStopOrder(stopOrder *pb.PostStopOrderRequest) (StopOrderID, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := sos.client.PostStopOrder(ctx, stopOrder)
	if err != nil {
		return "", err
	}

	return StopOrderID(res.StopOrderId), nil
}

func (sos StopOrdersService) GetStopOrders(accountID AccountID) ([]*pb.StopOrder, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := sos.client.GetStopOrders(ctx, &pb.GetStopOrdersRequest{
		AccountId: string(accountID),
	})
	if err != nil {
		return nil, err
	}

	return res.StopOrders, nil
}

func (sos StopOrdersService) CancelStopOrder(accountID AccountID, stopOrderID StopOrderID) (*timestamp.Timestamp, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := sos.client.CancelStopOrder(ctx, &pb.CancelStopOrderRequest{
		AccountId: string(accountID),
		StopOrderId: string(stopOrderID),
	})
	if err != nil {
		return nil, err
	}

	return res.Time, nil
}
