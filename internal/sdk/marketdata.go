package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// TODO: implementation
type MarketDataInterface interface {
	// Метод запроса исторических свечей по инструменту.
	GetCandles(figi Figi, from, to *timestamp.Timestamp, interval pb.CandleInterval) ([]*pb.HistoricCandle, error)
	// Метод запроса последних цен по инструментам.
	GetLastPrices(figi Figi) ([]*pb.LastPrice, error)
	// Метод получения стакана по инструменту.
	GetOrderBook(figi Figi, depth int) (*pb.OrderBook, error)
	// Метод запроса статуса торгов по инструментам.
	GetTradingStatus(figi Figi) (*pb.TradingStatus, error)
	// Метод запроса последних обезличенных сделок по инструменту.
	GetLastTrades(figi Figi, from, to *timestamp.Timestamp) ([]*pb.Trade, error)
}

type MarketDataService struct {
	client *pb.MarketDataServiceClient
}

func NewMarketDataService() *MarketDataService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewMarketDataServiceClient(conn)
	return &MarketDataService{client: &client}
}
