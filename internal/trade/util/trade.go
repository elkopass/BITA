// Package util stores some small and useful converters and formatters.
package util

import (
	"errors"
	pb "github.com/elkopass/BITA/internal/proto"
)

func GetVolumeAndLiquidity(candles []*pb.HistoricCandle) (int64, int64) {
	// Liquidity = (Q*V)/t
	var Q int64 = 0
	var V int64 = 0

	for _, candle := range candles {
		Q += candle.Volume
		V = (V + (candle.Open.Units+candle.Close.Units)/2) / 2
	}

	return Q, Q * V / 3600
}

func CalculateFairSellPrice(orderBook pb.GetOrderBookResponse) (*pb.Quotation, error) {
	if len(orderBook.Asks) == 0 {
		return nil, errors.New("no asks available")
	}

	return orderBook.Asks[0].Price, nil
}

func CalculateFairBuyPrice(orderBook pb.GetOrderBookResponse) (*pb.Quotation, error) {
	if len(orderBook.Bids) == 0 {
		return nil, errors.New("no bids available")
	}

	return orderBook.Bids[0].Price, nil
}
