// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type SandboxInterface interface {
	// The method of registering an account in the sandbox.
	OpenSandboxAccount() (string, error)
	// The method of getting accounts in the sandbox.
	GetSandboxAccounts() ([]*pb.Account, error)
	// The method of closing an account in the sandbox.
	CloseSandboxAccount(accountID string) error
	// The method of placing a trade order in the sandbox.
	PostSandboxOrder(order *pb.PostOrderRequest) (*pb.PostOrderResponse, error)
	// Method for getting a list of active applications for an account in the sandbox.
	GetSandboxOrders(accountID string) ([]*pb.OrderState, error)
	// Method for getting a list of active orders for an account in the sandbox.
	CancelSandboxOrder(accountID string, orderID string) (*timestamp.Timestamp, error)
	// The method of obtaining the order status in the sandbox.
	GetSandboxOrderState(accountID string, orderID string) (*pb.OrderState, error)
	// The method of obtaining positions on the virtual sandbox account.
	GetSandboxPositions(accountID string) (*pb.PositionsResponse, error)
	// The method of receiving operations in the sandbox by account number.
	GetSandboxOperations(filter *pb.OperationsRequest) ([]*pb.Operation, error)
	// The method of getting a portfolio in the sandbox.
	GetSandboxPortfolio(accountID string) (*pb.PortfolioResponse, error)
	// The method of depositing funds in the sandbox.
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

	ss.incrementRequestsCounter("OpenSandboxAccount")
	res, err := ss.client.OpenSandboxAccount(ctx, &pb.OpenSandboxAccountRequest{})
	if err != nil {
		ss.incrementApiCallErrors("OpenSandboxAccount", err.Error())
		return "", err
	}

	return res.AccountId, nil
}

func (ss SandboxService) GetSandboxAccounts() ([]*pb.Account, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("GetSandboxAccounts")
	res, err := ss.client.GetSandboxAccounts(ctx, &pb.GetAccountsRequest{})
	if err != nil {
		ss.incrementApiCallErrors("GetSandboxAccounts", err.Error())
		return nil, err
	}

	return res.Accounts, nil
}

func (ss SandboxService) CloseSandboxAccount(accountID string) error {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("CloseSandboxAccount")
	_, err := ss.client.CloseSandboxAccount(ctx, &pb.CloseSandboxAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		ss.incrementApiCallErrors("CloseSandboxAccount", err.Error())
		return err
	}

	return nil
}

func (ss SandboxService) PostSandboxOrder(order *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("PostSandboxOrder")
	res, err := ss.client.PostSandboxOrder(ctx, order)
	if err != nil {
		ss.incrementApiCallErrors("PostSandboxOrder", err.Error())
		return nil, err
	}

	return res, nil
}

func (ss SandboxService) GetSandboxOrders(accountID string) ([]*pb.OrderState, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("GetSandboxOrders")
	res, err := ss.client.GetSandboxOrders(ctx, &pb.GetOrdersRequest{
		AccountId: accountID,
	})
	if err != nil {
		ss.incrementApiCallErrors("GetSandboxOrders", err.Error())
		return nil, err
	}

	return res.Orders, nil
}

func (ss SandboxService) CancelSandboxOrder(accountID string, orderID string) (*timestamp.Timestamp, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("CancelSandboxOrder")
	res, err := ss.client.CancelSandboxOrder(ctx, &pb.CancelOrderRequest{
		AccountId: accountID,
		OrderId:   orderID,
	})
	if err != nil {
		ss.incrementApiCallErrors("CancelSandboxOrder", err.Error())
		return nil, err
	}

	return res.Time, nil
}

func (ss SandboxService) GetSandboxOrderState(accountID string, orderID string) (*pb.OrderState, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("GetSandboxOrderState")
	res, err := ss.client.GetSandboxOrderState(ctx, &pb.GetOrderStateRequest{
		AccountId: accountID,
		OrderId:   orderID,
	})
	if err != nil {
		ss.incrementApiCallErrors("GetSandboxOrderState", err.Error())
		return nil, err
	}

	return res, nil
}

func (ss SandboxService) GetSandboxPositions(accountID string) (*pb.PositionsResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("GetSandboxPositions")
	res, err := ss.client.GetSandboxPositions(ctx, &pb.PositionsRequest{
		AccountId: accountID,
	})
	if err != nil {
		ss.incrementApiCallErrors("GetSandboxPositions", err.Error())
		return nil, err
	}

	return res, nil
}

func (ss SandboxService) GetSandboxOperations(filter *pb.OperationsRequest) ([]*pb.Operation, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("GetSandboxOperations")
	res, err := ss.client.GetSandboxOperations(ctx, filter)
	if err != nil {
		ss.incrementApiCallErrors("GetSandboxOperations", err.Error())
		return nil, err
	}

	return res.Operations, nil
}

func (ss SandboxService) GetSandboxPortfolio(accountID string) (*pb.PortfolioResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("GetSandboxPortfolio")
	res, err := ss.client.GetSandboxPortfolio(ctx, &pb.PortfolioRequest{
		AccountId: accountID,
	})
	if err != nil {
		ss.incrementApiCallErrors("GetSandboxPortfolio", err.Error())
		return nil, err
	}

	return res, nil
}

func (ss SandboxService) SandboxPayIn(accountID string, amount *pb.MoneyValue) (*pb.MoneyValue, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	ss.incrementRequestsCounter("SandboxPayIn")
	res, err := ss.client.SandboxPayIn(ctx, &pb.SandboxPayInRequest{
		AccountId: accountID,
		Amount:    amount,
	})
	if err != nil {
		ss.incrementApiCallErrors("SandboxPayIn", err.Error())
		return nil, err
	}

	return res.Balance, nil
}

func (ss SandboxService) incrementRequestsCounter(method string) {
	metrics.ApiRequests.WithLabelValues(loggy.GetBotID(), "SandboxService", method).Inc()
}

func (ss SandboxService) incrementApiCallErrors(method string, error string) {
	metrics.ApiCallErrors.WithLabelValues(loggy.GetBotID(), "SandboxService", method, error).Inc()
}
