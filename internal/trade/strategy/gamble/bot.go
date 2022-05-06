package gamble

import (
	"context"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"go.uber.org/zap"
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
	tb.logger.Info("starting!")

	accountID, err := services.SandboxService.OpenSandboxAccount()
	if err != nil {
		tb.logger.Errorf("can not create account: %v", err)
	}
	tb.logger.Infof("created new account with ID %s", accountID)

	// replace logger
	tb.logger = tb.logger.With("account_id", accountID)

	wg := &sync.WaitGroup{}
	wg.Add(len(tb.config.Figi))

	for _, f := range tb.config.Figi {
		workerCtx, cancel := context.WithCancel(context.Background())

		w := NewTradeWorker(f, string(accountID))
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

	err = services.SandboxService.CloseSandboxAccount(accountID)
	if err != nil {
		tb.logger.Errorf("can't create account: %v", err)
	}
	tb.logger.Infof("account with ID %s closed successfully", accountID)

	return nil
}
