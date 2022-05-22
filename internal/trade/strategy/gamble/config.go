package gamble

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/sdk"
	"github.com/kelseyhightower/envconfig"
)

var services = sdk.NewServicePool()

type TradeConfig struct {
	LotsToBuy      int     `default:"1" split_words:"true"`
	StopLossCoef   float64 `default:"0.97" split_words:"true"`
	TakeProfitCoef float64 `default:"1.02" split_words:"true"`

	LongTrendToTrade  float64 `default:"0.05" split_words:"true"`
	ShortTrendToTrade float64 `default:"0.1" split_words:"true"`

	LongTrendIntervalSeconds  int `default:"86400" split_words:"true"`
	ShortTrendIntervalSeconds int `default:"3600" split_words:"true"`

	WorkerSleepDurationSeconds int64 `default:"30" split_words:"true"`
	SecondsToCancelOrder       int64 `default:"3600" split_words:"true"`
}

func NewTradeConfig() *TradeConfig {
	var c TradeConfig
	err := envconfig.Process("gamble_strategy", &c)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	return &c
}
