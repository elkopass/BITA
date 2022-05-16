// Package util stores some small and useful converters and formatters.
package util

import (
	"errors"
	pb "github.com/elkopass/BITA/internal/proto"
)

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
