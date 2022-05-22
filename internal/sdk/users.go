// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
)

type UsersServiceClient interface {
	// The method of receiving user accounts.
	GetAccounts() ([]*pb.Account, error)
	// Calculation of margin indicators on the account.
	GetMarginAttributes(accountID string) (*pb.GetMarginAttributesResponse, error)
	// Request for the user's tariff.
	GetUserTariff() (*pb.GetUserTariffResponse, error)
	// The method of obtaining information about the user.
	GetInfo() (*pb.GetInfoResponse, error)
}

type UsersService struct {
	client pb.UsersServiceClient
}

func NewUsersService() *UsersService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewUsersServiceClient(conn)
	return &UsersService{client: client}
}

func (us UsersService) GetAccounts() ([]*pb.Account, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	us.incrementRequestsCounter("GetAccounts")
	res, err := us.client.GetAccounts(ctx, &pb.GetAccountsRequest{})
	if err != nil {
		us.incrementApiCallErrors("GetAccounts", err.Error())
		return nil, err
	}

	return res.Accounts, nil
}

func (us UsersService) GetMarginAttributes(accountID string) (*pb.GetMarginAttributesResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	us.incrementRequestsCounter("GetMarginAttributes")
	res, err := us.client.GetMarginAttributes(ctx, &pb.GetMarginAttributesRequest{
		AccountId: accountID,
	})
	if err != nil {
		us.incrementApiCallErrors("GetMarginAttributes", err.Error())
		return nil, err
	}

	return res, nil
}

func (us UsersService) GetUserTariff() (*pb.GetUserTariffResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	us.incrementRequestsCounter("GetUserTariff")
	res, err := us.client.GetUserTariff(ctx, &pb.GetUserTariffRequest{})
	if err != nil {
		us.incrementApiCallErrors("GetUserTariff", err.Error())
		return nil, err
	}

	return res, nil
}

func (us UsersService) GetInfo() (*pb.GetInfoResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	us.incrementRequestsCounter("GetInfo")
	res, err := us.client.GetInfo(ctx, &pb.GetInfoRequest{})
	if err != nil {
		us.incrementApiCallErrors("GetInfo", err.Error())
		return nil, err
	}

	return res, nil
}

func (us UsersService) incrementRequestsCounter(method string) {
	metrics.ApiRequests.WithLabelValues(loggy.GetBotID(), "UsersService", method).Inc()
}

func (us UsersService) incrementApiCallErrors(method string, error string) {
	metrics.ApiCallErrors.WithLabelValues(loggy.GetBotID(), "UsersService", method, error).Inc()
}
