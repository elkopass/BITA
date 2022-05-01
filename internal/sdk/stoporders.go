package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// TODO: implementation
type StopOrdersInterface interface {
	// Метод выставления стоп-заявки.
	PostStopOrder(stopOrder *pb.PostStopOrderRequest) (StopOrderID, error)
	// Метод получения списка активных стоп заявок по счёту.
	GetStopOrders(accountID AccountID) ([]*pb.StopOrder, error)
	// Метод отмены стоп-заявки.
	CancelStopOrder(accountID AccountID, stopOrderID StopOrderID) (*timestamp.Timestamp, error)
}

type StopOrdersService struct {
	client *pb.StopOrdersServiceClient
}

func NewStopOrdersService() *StopOrdersService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewStopOrdersServiceClient(conn)
	return &StopOrdersService{client: &client}
}
