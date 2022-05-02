package config

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/kelseyhightower/envconfig"
	"time"
)

const (
	ApiURL         = "invest-public-api.tinkoff.ru:443"
	AppName        = "elkopass.TinkoffInvestRobotContest"
	DefaultTimeout = 30 * time.Second
)

type tradeBotConfig struct {
	IsSandbox  bool   `required:"true" split_words:"true"`
	Token      string `required:"true"`
	AccountID  string `split_words:"true"` // required in non-sandbox mode
	Env        string `default:"UNKNOWN"`
	Strategy   string `default:"gamble"`
	SellOnExit string `default:"false" split_words:"true"`
}

var (
	TradeBotConfig = func() tradeBotConfig {
		var config tradeBotConfig
		err := envconfig.Process("tradebot", &config)
		if err != nil {
			loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
		}

		return config
	}
)
