package gamble

import (
	"context"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/sdk"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

var services = sdk.NewServicePool()

type TradeStrategyConfig struct {
	Figi           []string `required:"true"`
	AmountToBuy    int      `split_words:"true"`
	TakeProfitCoef float64  `split_words:"true"`
	StopLossCoef   float64  `split_words:"true"`
	TrendToTrade   float64  `split_words:"true"`
}

func NewTradeConfig() *TradeStrategyConfig {
	var c TradeStrategyConfig
	err := envconfig.Process("gamble_strategy", &c)
	if err != nil {
		loggy.GetLogger().Sugar().Fatalf("failed to process config: %v", err)
	}

	return &c
}

type TradeWorker struct {
	ID      string
	Figi    string
	orderID string

	logger   *zap.SugaredLogger
}

func NewTradeWorker(figi, accountID string) *TradeWorker {
	id := strings.Split(uuid.New().String(), "-")[0]

	return &TradeWorker{
		ID:   id,
		Figi: figi,
		logger: loggy.GetLogger().Sugar().
			With("bot_id", loggy.GetBotID()).
			With("account_id", accountID).
			With("worker_id", id).
			With("figi", figi),
	}
}

func (tw TradeWorker) Run(ctx context.Context, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	tw.logger.Debug("start trading...")

	for {
		select {
		case <-time.After(3 * time.Second):
			orderBook, err := services.MarketDataService.GetOrderBook(sdk.Figi(tw.Figi), 1)
			if err != nil {
				tw.logger.Errorf("error getting orderBook: %v", err)
				continue // just ignoring it
			}

			tw.logger.Infof("last price: %d.%d", orderBook.LastPrice.Units, orderBook.LastPrice.Nano)
		case <-ctx.Done():
			tw.logger.Info("worker stopped!")
			return nil
		}
	}
}

type TradeBot struct {
	config      TradeStrategyConfig
	cancelFuncs []context.CancelFunc
	logger      *zap.SugaredLogger
}

func NewTradeBot() *TradeBot {
	return &TradeBot{
		config:   *NewTradeConfig(),
		logger:   loggy.GetLogger().Sugar().With("bot_id", loggy.GetBotID()),
	}
}

func (tb TradeBot) Run(ctx context.Context) (err error) {
	tb.logger.Info("starting in sandbox!")

	accountID, err := services.SandboxService.OpenSandboxAccount()
	if err != nil {
		tb.logger.Errorf("can not create account: %v", err)
	}
	tb.logger.Infof("created new account with ID %s", accountID)

	// replace logger
	tb.logger = tb.logger.With("account_id", accountID)

	res, err := services.SandboxService.SandboxPayIn(accountID, &pb.MoneyValue{
		Currency: "RUB",
		Units:    1000,
	})
	if err != nil {
		tb.logger.With("account_id", accountID).Errorf("can not pay in: %v", err)
	}
	tb.logger.Infof("account successfully replenished with %d.%d %s", res.Units, res.Nano, res.Currency)

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
