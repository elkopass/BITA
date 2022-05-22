package crumble

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/sdk"
	"github.com/kelseyhightower/envconfig"
)

var services = sdk.NewServicePool()

type TradeConfig struct {
	LotsToBuy      int     `default:"1" split_words:"true"`
	StopLossCoef   float64 `default:"0.95" split_words:"true"`
	TakeProfitCoef float64 `default:"1.05" split_words:"true"`

	ShortWindow          int `default:"25" split_words:"true"`
	LongWindow           int `default:"50" split_words:"true"`
	CandlesIntervalHours int `default:"144" split_words:"true"`

	WorkerSleepDurationSeconds int64 `default:"30" split_words:"true"`
	SecondsToCancelOrder       int64 `default:"3600" split_words:"true"`
}

func NewTradeConfig() *TradeConfig {
	var c TradeConfig
	err := envconfig.Process("crumble_strategy", &c)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	return &c
}
