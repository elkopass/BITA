package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type MarketDataInterface interface {
	// Метод запроса исторических свечей по инструменту.
	GetCandles(figi Figi, from, to *timestamp.Timestamp, interval pb.CandleInterval) ([]*pb.HistoricCandle, error)
	// Метод запроса последних цен по инструментам.
	GetLastPrices(figi []string) ([]*pb.LastPrice, error)
	// Метод получения стакана по инструменту.
	GetOrderBook(figi Figi, depth int) (*pb.GetOrderBookResponse, error)
	// Метод запроса статуса торгов по инструментам.
	GetTradingStatus(figi Figi) (*pb.GetTradingStatusResponse, error)
	// Метод запроса последних обезличенных сделок по инструменту.
	GetLastTrades(figi Figi, from, to *timestamp.Timestamp) ([]*pb.Trade, error)
}

type MarketDataService struct {
	client pb.MarketDataServiceClient
}

func NewMarketDataService() *MarketDataService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewMarketDataServiceClient(conn)
	return &MarketDataService{client: client}
}

func (mds MarketDataService) GetCandles(figi Figi, from, to *timestamp.Timestamp, interval pb.CandleInterval) ([]*pb.HistoricCandle, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := mds.client.GetCandles(ctx, &pb.GetCandlesRequest{
		Figi: string(figi),
		From: from,
		To: to,
		Interval: interval,
	})
	if err != nil {
		return nil, err
	}

	return res.Candles, nil
}

func (mds MarketDataService) GetLastPrices(figi []string) ([]*pb.LastPrice, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := mds.client.GetLastPrices(ctx, &pb.GetLastPricesRequest{
		Figi: figi,
	})
	if err != nil {
		return nil, err
	}

	return res.LastPrices, nil
}

func (mds MarketDataService) GetOrderBook(figi Figi, depth int) (*pb.GetOrderBookResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := mds.client.GetOrderBook(ctx, &pb.GetOrderBookRequest{
		Figi: string(figi),
		Depth: int32(depth),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (mds MarketDataService) GetTradingStatus(figi Figi) (*pb.GetTradingStatusResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := mds.client.GetTradingStatus(ctx, &pb.GetTradingStatusRequest{
		Figi: string(figi),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (mds MarketDataService) GetLastTrades(figi Figi, from, to *timestamp.Timestamp) ([]*pb.Trade, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := mds.client.GetLastTrades(ctx, &pb.GetLastTradesRequest{
		Figi: string(figi),
		From: from,
		To: to,
	})
	if err != nil {
		return nil, err
	}

	return res.Trades, nil
}
