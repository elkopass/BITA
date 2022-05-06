package gamble

import (
	"fmt"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	"strconv"
)

func candlesToTimeSeries(candles []*pb.HistoricCandle) *techan.TimeSeries {
	var techanCandles []*techan.Candle
	for i, c := range candles {
		if i == len(candles)-2 {
			break
		}

		tc := &techan.Candle{
			Period: techan.TimePeriod{
				Start: c.Time.AsTime(),
				End:   candles[i+1].Time.AsTime(),
			},
			Volume:     big.NewFromInt(int(c.Volume)),
			OpenPrice:  big.NewFromString(fmt.Sprintf("%d.%d", c.Close.Units, c.Close.Nano)),
			ClosePrice: big.NewFromString(fmt.Sprintf("%d.%d", c.Close.Units, c.Close.Nano)),
			MaxPrice:   big.NewFromString(fmt.Sprintf("%d.%d", c.High.Units, c.High.Nano)),
			MinPrice:   big.NewFromString(fmt.Sprintf("%d.%d", c.Low.Units, c.Low.Nano)),
		}
		techanCandles = append(techanCandles, tc)
	}

	return &techan.TimeSeries{Candles: techanCandles}
}

//func getFormattedCandles(candles []*pb.HistoricCandle) string {
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

func getFormattedPositions(positions []*pb.PortfolioPosition) string {
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

func QuotationToFloat(q pb.Quotation) float64 {
	if q.Nano < 0 {
		q.Nano = -q.Nano
	}
	p, _ := strconv.ParseFloat(fmt.Sprintf("%d.%d", q.Units, q.Nano), 64)

	return p
}

func MoneyValueToFloat(q pb.MoneyValue) float64 {
	if q.Nano < 0 {
		q.Nano = -q.Nano
	}
	p, _ := strconv.ParseFloat(fmt.Sprintf("%d.%d", q.Units, q.Nano), 64)

	return p
}