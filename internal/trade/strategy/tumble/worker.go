package tumble

import (
	"context"
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

type TradeWorker struct {
	ID        string
	Figi      string
	accountID string

	logger *zap.SugaredLogger
	config TradeConfig
}

func NewTradeWorker(figi, accountID string) *TradeWorker {
	id := strings.Split(uuid.New().String(), "-")[0]

	return &TradeWorker{
		ID:        id,
		Figi:      figi,
		accountID: accountID,
		config:    *NewTradeConfig(),
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
		case <-time.After(time.Duration(tw.config.WorkerSleepDurationSeconds) * time.Second):
			// TODO: implement trade logic
			tw.logger.Error("not implemented")
		case <-ctx.Done():
			// TODO: implement sell logic on interrupt
			tw.logger.Info("worker stopped!")

			return nil
		}
	}
}
