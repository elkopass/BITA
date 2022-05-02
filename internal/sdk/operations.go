package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// TODO: implementation
type OperationsInterface interface {
	// Метод получения списка операций по счёту.
	GetOperations(accountID AccountID, from, to *timestamp.Timestamp, state pb.OperationState, figi Figi) ([]*pb.Operation, error)
	// Метод получения портфеля по счёту.
	GetPortfolio(accountID AccountID) (*pb.PortfolioResponse, error)
	// Метод получения списка позиций по счёту.
	GetPositions(accountID AccountID) (*pb.PositionsResponse, error)
	// Метод получения доступного остатка для вывода средств.
	GetWithdrawLimits(accountID AccountID) (*pb.WithdrawLimitsResponse, error)
}

type OperationsService struct {
	client *pb.OperationsServiceClient
}

func NewOperationsService() *OperationsService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewOperationsServiceClient(conn)
	return &OperationsService{client: &client}
}
