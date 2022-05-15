package tumble

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/sdk"
	"github.com/kelseyhightower/envconfig"
)

var services = sdk.NewServicePool()

type TradeConfig struct {
	LotsToBuy int `default:"1" split_words:"true"`

	WorkerSleepDurationSeconds int `default:"30" split_words:"true"`
}

func NewTradeConfig() *TradeConfig {
	var c TradeConfig
	err := envconfig.Process("tumble_strategy", &c)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	return &c
}
