package gamble

import (
	"fmt"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
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

func getFormattedCandles(candles []*pb.HistoricCandle) string {
	formattedCandles := ""
	for _, candle := range candles {
		formattedCandles += fmt.Sprintf(
			"(%s) %d.%d",
			candle.Time.AsTime().String(),
			candle.Close.Units,
			candle.Close.Nano,
		)
	}

	return formattedCandles
}