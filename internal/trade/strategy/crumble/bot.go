package crumble

import (
	"context"
	"fmt"
	"github.com/elkopass/BITA/internal/config"
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/sdk"
	"go.uber.org/zap"
	"strings"
	"sync"
)

type TradeBot struct {
	config      TradeConfig
	cancelFuncs []context.CancelFunc
	logger      *zap.SugaredLogger
}

func NewTradeBot() *TradeBot {
	return &TradeBot{
		config: *NewTradeConfig(),
		logger: loggy.GetLogger().Sugar().With("bot_id", loggy.GetBotID()),
	}
}

func (tb TradeBot) Run(ctx context.Context) (err error) {
	tb.logger.Infof("starting with %s strategy and sdk v%s", config.TradeBotConfig().Strategy, sdk.Version)

	accountID := config.TradeBotConfig().AccountID
	if config.TradeBotConfig().IsSandbox {
		accountID, err = services.SandboxService.OpenSandboxAccount()
		if err != nil {
			return fmt.Errorf("can not create account: %v", err)
		}
		tb.logger.Infof("created new account with ID %s", accountID)
	} else {
		info, err := services.UsersService.GetInfo()
		if err != nil {
			return fmt.Errorf("can not get user info: %v", err)
		}
		tb.logger.Infof("user tariff: %s, qualified for work with %s",
			info.Tariff, strings.Join(info.QualifiedForWorkWith, ","))
	}

	// replace logger
	tb.logger = tb.logger.With("account_id", accountID)

	figi := config.TradeBotConfig().Figi
	wg := &sync.WaitGroup{}
	wg.Add(len(figi))

	for _, f := range figi {
		workerCtx, cancel := context.WithCancel(context.Background())

		w := NewTradeWorker(f, accountID)
		tb.cancelFuncs = append(tb.cancelFuncs, cancel)

		go func() {
			err = w.Run(workerCtx, wg)
			if err != nil {
				tb.logger.Errorf("worker finished with error: %v", err)
			}
		}()
	}

	<-ctx.Done()

	for _, cancel := range tb.cancelFuncs {
		cancel()
	}

	wg.Wait()

	if config.TradeBotConfig().IsSandbox {
		err = services.SandboxService.CloseSandboxAccount(accountID)
		if err != nil {
			tb.logger.Errorf("can't close an account: %v", err)
		}
		tb.logger.Infof("account with ID %s closed successfully", accountID)
	}

	return nil
}
