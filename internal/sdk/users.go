package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
)

type UsersServiceClient interface {
	// Метод получения счетов пользователя.
	GetAccounts() ([]*pb.Account, error)
	// Расчёт маржинальных показателей по счёту.
	GetMarginAttributes(accountID AccountID) (*pb.GetMarginAttributesResponse, error)
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

	res, err := us.client.GetAccounts(ctx, &pb.GetAccountsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Accounts, nil
}

func (us UsersService) GetMarginAttributes(accountID AccountID) (*pb.GetMarginAttributesResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := us.client.GetMarginAttributes(ctx, &pb.GetMarginAttributesRequest{
		AccountId: string(accountID),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (us UsersService) GetUserTariff() (*pb.GetUserTariffResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := us.client.GetUserTariff(ctx, &pb.GetUserTariffRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (us UsersService) GetInfo() (*pb.GetInfoResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := us.client.GetInfo(ctx, &pb.GetInfoRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
