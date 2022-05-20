// Package config stores global trade-bot configuration.
package config

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/kelseyhightower/envconfig"
	"time"
)

const (
	ApiURL  = "invest-public-api.tinkoff.ru:443"
	AppName = "elkopass.BITA"

	DefaultRequestTimeout = 30 * time.Second

	CircuitBreakerMaxFailures = 5
	CircuitBreakerRefreshTime = 30 * time.Minute
)

type tradeBotConfig struct {
	Figi []string `required:"true"`

	IsSandbox    bool   `required:"true" split_words:"true"`
	Token        string `required:"true"`
	AccountID    string `split_words:"true"` // required in non-sandbox mode
	Env          string `default:"UNSPECIFIED"`
	LogLevel     string `default:"DEBUG" split_words:"true"`
	Strategy     string `default:"gamble"`
	SellOnExit   bool   `default:"false" split_words:"true"`
	TimeToCancel int64  `default:"3600" split_words:"true"`
}

type metricsConfig struct {
	Enabled  bool   `default:"true" split_words:"true"`
	Addr     string `default:":8080" split_words:"true"`
	Endpoint string `default:"/metrics" split_words:"true"`
}

var (
	// TradeBotConfig returns relevant global configuration.
	TradeBotConfig = func() tradeBotConfig {
		var config tradeBotConfig
		err := envconfig.Process("tradebot", &config)
		if err != nil {
			loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
		}

		return config
	}

	// TradeBotConfig returns config for Prometheus exporter.
	MetricsConfig = func() metricsConfig {
		var config metricsConfig
		err := envconfig.Process("metrics", &config)
		if err != nil {
			loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
		}

		return config
	}
)
