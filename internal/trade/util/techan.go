// Package util stores some small and useful converters and formatters.
package util

import (
	"fmt"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

// CandlesToTimeSeries converts historic candles from Invest API to techan.TimeSeries.
func CandlesToTimeSeries(candles []*pb.HistoricCandle) *techan.TimeSeries {
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
			OpenPrice:  big.NewFromString(fmt.Sprintf("%d.%d", c.Close.Units, abs(c.Close.Nano))),
			ClosePrice: big.NewFromString(fmt.Sprintf("%d.%d", c.Close.Units, abs(c.Close.Nano))),
			MaxPrice:   big.NewFromString(fmt.Sprintf("%d.%d", c.High.Units, abs(c.High.Nano))),
			MinPrice:   big.NewFromString(fmt.Sprintf("%d.%d", c.Low.Units, abs(c.Low.Nano))),
		}
		techanCandles = append(techanCandles, tc)
	}

	return &techan.TimeSeries{Candles: techanCandles}
}

func abs(x int32) int32 {
	return absDiff(x, 0)
}

func absDiff(x, y int32) int32 {
	if x < y {
		return y - x
	}
	return x - y
}
