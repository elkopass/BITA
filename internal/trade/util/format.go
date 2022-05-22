// Package util stores some small and useful converters and formatters.
package util

import (
	"fmt"
	pb "github.com/elkopass/BITA/internal/proto"
)

// GetFormattedCandles returns historic candles in a pretty-formatted string to print.
//func GetFormattedCandles(candles []*pb.HistoricCandle) string {
//	formattedCandles := ""
//	for _, c := range candles {
//		formattedCandles += fmt.Sprintf(
//			"(%s) %d.%d",
//			c.Time.AsTime().String(),
//			c.Close.Units,
//			c.Close.Nano,
//		)
//	}
//
//	return formattedCandles
//}

// GetFormattedPositions returns portfolio positions in a pretty-formatted string to print.
func GetFormattedPositions(positions []*pb.PortfolioPosition) string {
	formattedPositions := ""
	for _, p := range positions {
		formattedPositions += fmt.Sprintf(
			"%s (%s)",
			p.Figi,
			p.InstrumentType,
		)
		if p.Quantity != nil {
			formattedPositions += fmt.Sprintf(
				", %d.%d quantity",
				p.Quantity.Units,
				p.Quantity.Nano,
			)
		}
		if p.CurrentPrice != nil {
			formattedPositions += fmt.Sprintf(
				", %d.%d price",
				p.CurrentPrice.Units,
				p.CurrentPrice.Nano,
			)
		}
		if p.ExpectedYield != nil {
			formattedPositions += fmt.Sprintf(
				", %d.%d yield",
				p.ExpectedYield.Units,
				p.ExpectedYield.Nano,
			)
		}

		formattedPositions += ";"
	}

	return formattedPositions
}

// GetFormattedOrderBook returns order book in a pretty-formatted string to print.
func GetFormattedOrderBook(orderBook *pb.OrderBook) string {
	formattedOrderBook := "bids: "
	for _, b := range orderBook.Bids {
		price := QuotationToFloat(*b.Price)
		formattedOrderBook += fmt.Sprintf("%f (%d) ", price, b.Quantity)
	}

	formattedOrderBook += " | asks: "
	for _, a := range orderBook.Asks {
		price := QuotationToFloat(*a.Price)
		formattedOrderBook += fmt.Sprintf("%f (%d) ", price, a.Quantity)
	}

	return formattedOrderBook
}
