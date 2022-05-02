package config

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/kelseyhightower/envconfig"
	"time"
)

const (
	ApiURL  = "invest-public-api.tinkoff.ru:443"
	AppName = "elkopass.TinkoffInvestRobotContest"

	DefaultRequestTimeout   = 30 * time.Second
	// GracefulShutdownTimeout = 60 * time.Second
)

type tradeBotConfig struct {
	IsSandbox  bool   `required:"true" split_words:"true"`
	Token      string `required:"true"`
	AccountID  string `split_words:"true"` // required in non-sandbox mode
	Env        string `default:"UNKNOWN"`
	Strategy   string `default:"gamble"`
	SellOnExit string `default:"false" split_words:"true"`
}

type metricsConfig struct {
	Enabled  bool   `default:"true" split_words:"true"`
	Addr     string `default:":8080" split_words:"true"`
	Endpoint string `default:"/metrics" split_words:"true"`
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

	MetricsConfig = func() metricsConfig {
		var config metricsConfig
		err := envconfig.Process("metrics", &config)
		if err != nil {
			loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
		}

		return config
	}
)
