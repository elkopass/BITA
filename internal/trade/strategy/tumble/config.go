package tumble

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/sdk"
	"github.com/kelseyhightower/envconfig"
)

var services = sdk.NewServicePool()

type TradeConfig struct {
	LotsToBuy int `default:"1" split_words:"true"`

	AsksBidsRatio float64 `default:"1.5" split_words:"true"`
	BidsAsksRatio float64 `default:"1.5" split_words:"true"`

	OrderBookDepth        int `default:"10" split_words:"true"`
	OrderBookFairAskDepth int `default:"5" split_words:"true"`
	OrderBookFairBidDepth int `default:"5" split_words:"true"`
}

func NewTradeConfig() *TradeConfig {
	var c TradeConfig
	err := envconfig.Process("tumble_strategy", &c)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	if c.OrderBookFairBidDepth > c.OrderBookDepth {
		loggy.GetLogger().Sugar().Fatal("fair bid depth must be lower than order book depth")
	}
	if c.OrderBookFairAskDepth > c.OrderBookDepth {
		loggy.GetLogger().Sugar().Fatal("fair ask depth must be lower than order book depth")
	}

	return &c
}
