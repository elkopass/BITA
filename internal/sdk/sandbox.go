package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// TODO: implementation
type SandboxInterface interface {
	// Метод регистрации счёта в песочнице.
	OpenSandboxAccount() (AccountID, error)
	// Метод получения счетов в песочнице.
	GetSandboxAccounts() ([]*pb.Account, error)
	// Метод закрытия счёта в песочнице.
	CloseSandboxAccount(accountID AccountID) error
	// Метод выставления торгового поручения в песочнице.
	PostSandboxOrder(order *Order) (*PostOrder, error)
	// Метод получения списка активных заявок по счёту в песочнице.
	GetSandboxOrders(accountID AccountID) ([]*OrderState, error)
	// Метод отмены торгового поручения в песочнице.
	CancelSandboxOrder(accountID AccountID, orderID OrderID) (*timestamp.Timestamp, error)
	// Метод получения статуса заявки в песочнице.
	GetSandboxOrderState(accountID AccountID, orderID OrderID) (*OrderState, error)
	// Метод получения позиций по виртуальному счёту песочницы.
	GetSandboxPositions(accountID AccountID) (*Positions, error)
	// Метод получения операций в песочнице по номеру счёта.
	GetSandboxOperations(filter *OperationsSearchFilters) ([]*Operation, error)
	// Метод получения портфолио в песочнице.
	GetSandboxPortfolio(accountID AccountID) (*Portfolio, error)
	// Метод пополнения счёта в песочнице.
	SandboxPayIn(accountID AccountID, amount *pb.MoneyValue) (*pb.MoneyValue, error)
}

type SandboxService struct {
	client *pb.SandboxServiceClient
}

func NewSandboxService() *SandboxService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewSandboxServiceClient(conn)
	return &SandboxService{client: &client}
}
