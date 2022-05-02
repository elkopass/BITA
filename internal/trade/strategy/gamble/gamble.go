package gamble

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"time"
)

type TraderStrategyConfig struct {
	Figi           []string `required:"true"`
	AmountToBuy    int      `split_words:"true"`
	TakeProfitCoef float64  `split_words:"true"`
	StopLossCoef   float64  `split_words:"true"`
	TrendToTrade   float64  `split_words:"true"`
}

type TraderState struct {
	currentPurchasePrice  *float64
	currentPurchaseAmount *int
	lastPurchaseTime      *time.Time
	lastSellTime          *time.Time
}

type TraderBot struct {
	config TraderStrategyConfig
	state  TraderState
	logger *zap.SugaredLogger
}

func NewTraderConfig() *TraderStrategyConfig {
	var c TraderStrategyConfig
	err := envconfig.Process("gamble_strategy", &c)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	return &c
}

func NewTraderState() *TraderState {
	return &TraderState{}
}

func NewTraderBot() *TraderBot {
	return &TraderBot{
		config: *NewTraderConfig(),
		state:  *NewTraderState(),
		logger: loggy.GetLogger().Sugar(),
	}
}

func (tb TraderBot) Run() {
	tb.logger.Info("Starting!")
}
