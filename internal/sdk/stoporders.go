// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type StopOrdersInterface interface {
	// The method of placing a stop order.
	PostStopOrder(stopOrder *pb.PostStopOrderRequest) (string, error)
	// Method for getting a list of active stop orders on the account.
	GetStopOrders(accountID string) ([]*pb.StopOrder, error)
	// The method of canceling the stop order.
	CancelStopOrder(accountID string, stopOrderID string) (*timestamp.Timestamp, error)
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

func (sos StopOrdersService) PostStopOrder(stopOrder *pb.PostStopOrderRequest) (string, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	sos.incrementRequestsCounter("PostStopOrder")
	res, err := sos.client.PostStopOrder(ctx, stopOrder)
	if err != nil {
		sos.incrementApiCallErrors("PostStopOrder", err.Error())
		return "", err
	}

	return res.StopOrderId, nil
}

func (sos StopOrdersService) GetStopOrders(accountID string) ([]*pb.StopOrder, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	sos.incrementRequestsCounter("GetStopOrders")
	res, err := sos.client.GetStopOrders(ctx, &pb.GetStopOrdersRequest{
		AccountId: accountID,
	})
	if err != nil {
		sos.incrementApiCallErrors("GetStopOrders", err.Error())
		return nil, err
	}

	return res.StopOrders, nil
}

func (sos StopOrdersService) CancelStopOrder(accountID string, stopOrderID string) (*timestamp.Timestamp, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	sos.incrementRequestsCounter("CancelStopOrder")
	res, err := sos.client.CancelStopOrder(ctx, &pb.CancelStopOrderRequest{
		AccountId:   accountID,
		StopOrderId: stopOrderID,
	})
	if err != nil {
		sos.incrementApiCallErrors("CancelStopOrder", err.Error())
		return nil, err
	}

	return res.Time, nil
}

func (sos StopOrdersService) incrementRequestsCounter(method string) {
	metrics.ApiRequests.WithLabelValues(loggy.GetBotID(), "StopOrdersService", method).Inc()
}

func (sos StopOrdersService) incrementApiCallErrors(method string, error string) {
	metrics.ApiCallErrors.WithLabelValues(loggy.GetBotID(), "StopOrdersService", method, error).Inc()
}
