package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type SandboxInterface interface {
	// Метод регистрации счёта в песочнице.
	OpenSandboxAccount() (string, error)
	// Метод получения счетов в песочнице.
	GetSandboxAccounts() ([]*pb.Account, error)
	// Метод закрытия счёта в песочнице.
	CloseSandboxAccount(accountID string) error
	// Метод выставления торгового поручения в песочнице.
	PostSandboxOrder(order *pb.PostOrderRequest) (*pb.PostOrderResponse, error)
	// Метод получения списка активных заявок по счёту в песочнице.
	GetSandboxOrders(accountID string) ([]*pb.OrderState, error)
	// Метод отмены торгового поручения в песочнице.
	CancelSandboxOrder(accountID string, orderID string) (*timestamp.Timestamp, error)
	// Метод получения статуса заявки в песочнице.
	GetSandboxOrderState(accountID string, orderID string) (*pb.OrderState, error)
	// Метод получения позиций по виртуальному счёту песочницы.
	GetSandboxPositions(accountID string) (*pb.PositionsResponse, error)
	// Метод получения операций в песочнице по номеру счёта.
	GetSandboxOperations(filter *pb.OperationsRequest) ([]*pb.Operation, error)
	// Метод получения портфолио в песочнице.
	GetSandboxPortfolio(accountID string) (*pb.PortfolioResponse, error)
	// Метод пополнения счёта в песочнице.
	SandboxPayIn(accountID string, amount *pb.MoneyValue) (*pb.MoneyValue, error)
}

type SandboxService struct {
	client pb.SandboxServiceClient
}

func NewSandboxService() *SandboxService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewSandboxServiceClient(conn)
	return &SandboxService{client: client}
}

func (ss SandboxService) OpenSandboxAccount() (string, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.OpenSandboxAccount(ctx, &pb.OpenSandboxAccountRequest{})
	if err != nil {
		return "", err
	}

	return string(res.AccountId), nil
}

func (ss SandboxService) GetSandboxAccounts() ([]*pb.Account, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.GetSandboxAccounts(ctx, &pb.GetAccountsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Accounts, nil
}

func (ss SandboxService) CloseSandboxAccount(accountID string) error {
	ctx, cancel := createRequestContext()
	defer cancel()

	_, err := ss.client.CloseSandboxAccount(ctx, &pb.CloseSandboxAccountRequest{
		AccountId: accountID,
	})
	return err
}

func (ss SandboxService) PostSandboxOrder(order *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.PostSandboxOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ss SandboxService) GetSandboxOrders(accountID string) ([]*pb.OrderState, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.GetSandboxOrders(ctx, &pb.GetOrdersRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, err
	}

	return res.Orders, nil
}

func (ss SandboxService) CancelSandboxOrder(accountID string, orderID string) (*timestamp.Timestamp, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.CancelSandboxOrder(ctx, &pb.CancelOrderRequest{
		AccountId: accountID,
		OrderId:   orderID,
	})
	if err != nil {
		return nil, err
	}

	return res.Time, nil
}

func (ss SandboxService) GetSandboxOrderState(accountID string, orderID string) (*pb.OrderState, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.GetSandboxOrderState(ctx, &pb.GetOrderStateRequest{
		AccountId: accountID,
		OrderId:   orderID,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ss SandboxService) GetSandboxPositions(accountID string) (*pb.PositionsResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.GetSandboxPositions(ctx, &pb.PositionsRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ss SandboxService) GetSandboxOperations(filter *pb.OperationsRequest) ([]*pb.Operation, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.GetSandboxOperations(ctx, filter)
	if err != nil {
		return nil, err
	}

	return res.Operations, nil
}

func (ss SandboxService) GetSandboxPortfolio(accountID string) (*pb.PortfolioResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.GetSandboxPortfolio(ctx, &pb.PortfolioRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ss SandboxService) SandboxPayIn(accountID string, amount *pb.MoneyValue) (*pb.MoneyValue, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := ss.client.SandboxPayIn(ctx, &pb.SandboxPayInRequest{
		AccountId: accountID,
		Amount:    amount,
	})
	if err != nil {
		return nil, err
	}

	return res.Balance, nil
}
