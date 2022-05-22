// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type MarketDataInterface interface {
	// The method of requesting historical candlesticks by instrument.
	GetCandles(figi string, from, to *timestamp.Timestamp, interval pb.CandleInterval) ([]*pb.HistoricCandle, error)
	// The method of requesting the latest prices for instruments.
	GetLastPrices(figi []string) ([]*pb.LastPrice, error)
	// The method of obtaining a glass by instrument.
	GetOrderBook(figi string, depth int) (*pb.GetOrderBookResponse, error)
	// The method of requesting the status of trading on instruments.
	GetTradingStatus(figi string) (*pb.GetTradingStatusResponse, error)
	// The method of requesting the latest depersonalized transactions on the instrument.
	GetLastTrades(figi string, from, to *timestamp.Timestamp) ([]*pb.Trade, error)
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

func (mds MarketDataService) GetCandles(figi string, from, to *timestamp.Timestamp, interval pb.CandleInterval) ([]*pb.HistoricCandle, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	mds.incrementRequestsCounter("GetCandles")
	res, err := mds.client.GetCandles(ctx, &pb.GetCandlesRequest{
		Figi:     figi,
		From:     from,
		To:       to,
		Interval: interval,
	})
	if err != nil {
		mds.incrementApiCallErrors("GetCandles", err.Error())
		return nil, err
	}

	return res.Candles, nil
}

func (mds MarketDataService) GetLastPrices(figi []string) ([]*pb.LastPrice, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	mds.incrementRequestsCounter("GetLastPrices")
	res, err := mds.client.GetLastPrices(ctx, &pb.GetLastPricesRequest{
		Figi: figi,
	})
	if err != nil {
		mds.incrementApiCallErrors("GetLastPrices", err.Error())
		return nil, err
	}

	return res.LastPrices, nil
}

func (mds MarketDataService) GetOrderBook(figi string, depth int) (*pb.GetOrderBookResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	mds.incrementRequestsCounter("GetOrderBook")
	res, err := mds.client.GetOrderBook(ctx, &pb.GetOrderBookRequest{
		Figi:  figi,
		Depth: int32(depth),
	})
	if err != nil {
		mds.incrementApiCallErrors("GetOrderBook", err.Error())
		return nil, err
	}

	return res, nil
}

func (mds MarketDataService) GetTradingStatus(figi string) (*pb.GetTradingStatusResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	mds.incrementRequestsCounter("GetTradingStatus")
	res, err := mds.client.GetTradingStatus(ctx, &pb.GetTradingStatusRequest{
		Figi: figi,
	})
	if err != nil {
		mds.incrementApiCallErrors("GetTradingStatus", err.Error())
		return nil, err
	}

	return res, nil
}

func (mds MarketDataService) GetLastTrades(figi string, from, to *timestamp.Timestamp) ([]*pb.Trade, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	mds.incrementRequestsCounter("GetLastTrades")
	res, err := mds.client.GetLastTrades(ctx, &pb.GetLastTradesRequest{
		Figi: figi,
		From: from,
		To:   to,
	})
	if err != nil {
		mds.incrementApiCallErrors("GetLastTrades", err.Error())
		return nil, err
	}

	return res.Trades, nil
}

func (mds MarketDataService) incrementRequestsCounter(method string) {
	metrics.ApiRequests.WithLabelValues(loggy.GetBotID(), "MarketDataService", method).Inc()
}

func (mds MarketDataService) incrementApiCallErrors(method string, error string) {
	metrics.ApiCallErrors.WithLabelValues(loggy.GetBotID(), "MarketDataService", method, error).Inc()
}
