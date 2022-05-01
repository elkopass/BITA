package config

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/kelseyhightower/envconfig"
	"time"
)

const (
	ApiURL  = "invest-public-api.tinkoff.ru:443"
	AppName = "elkopass.TinkoffInvestRobotContest"
	DefaultTimeout = 30*time.Second
)

type tradeBotConfig struct {
	Token string
	Env   string `default:"UNKNOWN"`
}

var (
	TradeBotConfig = func() tradeBotConfig {
		var tradeBotConfig tradeBotConfig
		err := envconfig.Process("tradebot", &tradeBotConfig)
		if err != nil {
			loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
		}

		return tradeBotConfig
	}
)

