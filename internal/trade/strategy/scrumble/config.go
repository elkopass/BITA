package scrumble

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

	MMAIntervalSeconds int `default:"518400" split_words:"true"`

	WorkerSleepDurationSeconds int `default:"10" split_words:"true"`
}

func NewTradeConfig() *TradeConfig {
	var c TradeConfig
	err := envconfig.Process("scrumble_strategy", &c)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	return &c
}
