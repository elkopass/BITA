// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type OperationsInterface interface {
	// Method for getting a list of account transactions.
	GetOperations(accountID string, from, to *timestamp.Timestamp, state pb.OperationState, figi string) ([]*pb.Operation, error)
	// The method of obtaining a portfolio by account.
	GetPortfolio(accountID string) (*pb.PortfolioResponse, error)
	// Method for getting a list of account positions.
	GetPositions(accountID string) (*pb.PositionsResponse, error)
	// The method of obtaining the available balance for withdrawal of funds.
	GetWithdrawLimits(accountID string) (*pb.WithdrawLimitsResponse, error)
}

type OperationsService struct {
	client pb.OperationsServiceClient
}

func NewOperationsService() *OperationsService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewOperationsServiceClient(conn)
	return &OperationsService{client: client}
}

func (os OperationsService) GetOperations(accountID string, from, to *timestamp.Timestamp, state pb.OperationState, figi string) ([]*pb.Operation, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	os.incrementRequestsCounter("GetOperations")
	res, err := os.client.GetOperations(ctx, &pb.OperationsRequest{
		AccountId: accountID,
		From:      from,
		To:        to,
		State:     state,
		Figi:      figi,
	})
	if err != nil {
		os.incrementApiCallErrors("GetOperations", err.Error())
		return nil, err
	}

	return res.Operations, nil
}

func (os OperationsService) GetPortfolio(accountID string) (*pb.PortfolioResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	os.incrementRequestsCounter("GetPortfolio")
	res, err := os.client.GetPortfolio(ctx, &pb.PortfolioRequest{
		AccountId: accountID,
	})
	if err != nil {
		os.incrementApiCallErrors("GetPortfolio", err.Error())
		return nil, err
	}

	return res, nil
}

func (os OperationsService) GetPositions(accountID string) (*pb.PositionsResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	os.incrementRequestsCounter("GetPositions")
	res, err := os.client.GetPositions(ctx, &pb.PositionsRequest{
		AccountId: accountID,
	})
	if err != nil {
		os.incrementApiCallErrors("GetPositions", err.Error())
		return nil, err
	}

	return res, nil
}

func (os OperationsService) GetWithdrawLimits(accountID string) (*pb.WithdrawLimitsResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	os.incrementRequestsCounter("GetWithdrawLimits")
	res, err := os.client.GetWithdrawLimits(ctx, &pb.WithdrawLimitsRequest{
		AccountId: accountID,
	})
	if err != nil {
		os.incrementApiCallErrors("GetWithdrawLimits", err.Error())
		return nil, err
	}

	return res, nil
}

func (os OperationsService) incrementRequestsCounter(method string) {
	metrics.ApiRequests.WithLabelValues(loggy.GetBotID(), "OperationsService", method).Inc()
}

func (os OperationsService) incrementApiCallErrors(method string, error string) {
	metrics.ApiCallErrors.WithLabelValues(loggy.GetBotID(), "OperationsService", method, error).Inc()
}
