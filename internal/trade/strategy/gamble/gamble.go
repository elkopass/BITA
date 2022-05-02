package gamble

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/sdk"
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

	sandboxService sdk.SandboxService
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

		sandboxService: *sdk.NewSandboxService(),
	}
}

func (tb TraderBot) Run() {
	tb.logger.Info("Starting!")
}

func (tb TraderBot) RunInSandbox() {
	tb.logger.Info("starting in sandbox!")

	accountID, err := tb.sandboxService.OpenSandboxAccount()
	if err != nil {
		tb.logger.Errorf("can not create account: %v", err)
	}
	tb.logger.Infof("created new account with ID %s", accountID)

	// replace logger
	tb.logger = tb.logger.With("account_id", accountID)

	res, err := tb.sandboxService.SandboxPayIn(accountID, &pb.MoneyValue{
		Currency: "RUB",
		Units:    1000,
	})
	if err != nil {
		tb.logger.With("account_id", accountID).Errorf("can not pay in: %v", err)
	}
	tb.logger.Infof("account successfully replenished with %d.%d %s", res.Units, res.Nano, res.Currency)

	err = tb.sandboxService.CloseSandboxAccount(accountID)
	if err != nil {
		tb.logger.Errorf("can't create account: %v", err)
	}
	tb.logger.Infof("account with ID %s closed successfully", accountID)
}
