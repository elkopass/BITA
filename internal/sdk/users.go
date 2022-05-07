package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/metrics"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
)

type UsersServiceClient interface {
	// Метод получения счетов пользователя.
	GetAccounts() ([]*pb.Account, error)
	// Расчёт маржинальных показателей по счёту.
	GetMarginAttributes(accountID string) (*pb.GetMarginAttributesResponse, error)
	// Запрос тарифа пользователя.
	GetUserTariff() (*pb.GetUserTariffResponse, error)
	// Метод получения информации о пользователе.
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
