package gamble

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/sdk"
	"github.com/kelseyhightower/envconfig"
)

var services = sdk.NewServicePool()

type TradeConfig struct {
	Figi []string `required:"true"`

	AmountToBuy    int     `default:"1" split_words:"true"`
	StopLossCoef   float64 `default:"0.99" split_words:"true"`
	TakeProfitCoef float64 `default:"1.01" split_words:"true"`

	LongTrendToTrade  float64 `default:"0.05" split_words:"true"`
	ShortTrendToTrade float64 `default:"0.1" split_words:"true"`

	LongTrendIntervalSeconds  int `default:"86400" split_words:"true"`
	ShortTrendIntervalSeconds int `default:"3600" split_words:"true"`

	WorkerSleepDurationSeconds int `default:"30" split_words:"true"`
}

func NewTradeConfig() *TradeConfig {
	var c TradeConfig
	err := envconfig.Process("gamble_strategy", &c)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	return &c
}
