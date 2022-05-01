package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
)

// TODO: implementation
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
	client *pb.UsersServiceClient
}

func NewUsersService() *UsersService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewUsersServiceClient(conn)
	return &UsersService{client: &client}
}
