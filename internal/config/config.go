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
)

type tradeBotConfig struct {
	Figi []string `split_words:"true"`

	IsSandbox  bool   `default:"true" split_words:"true"`
	Token      string `required:"true"`
	AccountID  string `split_words:"true"` // required in non-sandbox mode
	Env        string `default:"UNSPECIFIED"`
	LogLevel   string `default:"INFO" split_words:"true"`
	Strategy   string `default:"gamble"`
	SellOnExit bool   `default:"false" split_words:"true"`
}

type metricsConfig struct {
	Enabled  bool   `default:"true" split_words:"true"`
	Addr     string `default:":8080" split_words:"true"`
	Endpoint string `default:"/metrics" split_words:"true"`
}

type circuitBreakerConfig struct {
	MaxFailures        int `default:"5" split_words:"true"`
	RefreshTimeMinutes int `default:"60" split_words:"true"`
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

	// MetricsConfig returns config for Prometheus exporter.
	MetricsConfig = func() metricsConfig {
		var config metricsConfig
		err := envconfig.Process("metrics", &config)
		if err != nil {
			loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
		}

		return config
	}

	// CircuitBreakerConfig returns config for breaker.CircuitBreaker.
	CircuitBreakerConfig = func() circuitBreakerConfig {
		var config circuitBreakerConfig
		err := envconfig.Process("breaker", &config)
		if err != nil {
			loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
		}

		return config
	}
)
