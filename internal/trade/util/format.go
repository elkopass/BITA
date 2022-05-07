package util

import (
	"fmt"
	pb "github.com/elkopass/BITA/internal/proto"
)

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
