package config

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/kelseyhightower/envconfig"
)

const (
	ApiURL  = "invest-public-api.tinkoff.ru:443"
	AppName = "elkopass.TinkoffInvestRobotContest"
)

type TradeBotConfig struct {
	Token string
	Env   string `default:"UNKNOWN"`
}

func GetTradeBotConfig() TradeBotConfig {
	var tradeBotConfig TradeBotConfig
	err := envconfig.Process("tradebot", &tradeBotConfig)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	return tradeBotConfig
}
