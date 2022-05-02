package gamble

import (
	"context"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/sdk"
	"github.com/google/uuid"
	"github.com/sdcoffey/techan"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"sync"
	"time"
)

type TradeWorker struct {
	ID      string
	Figi    string
	orderID string

	logger *zap.SugaredLogger
	config TradeConfig
}

func NewTradeWorker(figi, accountID string) *TradeWorker {
	id := strings.Split(uuid.New().String(), "-")[0]

	return &TradeWorker{
		ID:     id,
		Figi:   figi,
		config: *NewTradeConfig(),
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
			status, err := services.MarketDataService.GetTradingStatus(sdk.Figi(tw.Figi))
			if err != nil {
				tw.logger.Errorf("error getting trading status: %v", err)
				continue // just ignoring it
			}

			tw.logger.Infof("trading status: %s", status.TradingStatus.String())
			if status.TradingStatus != pb.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING {
				continue // wait until trading is available
			}

			candles, err := services.MarketDataService.GetCandles(
				sdk.Figi(tw.Figi),
				timestamppb.New(time.Now().Add(-1*time.Hour)),
				timestamppb.Now(),
				pb.CandleInterval_CANDLE_INTERVAL_5_MIN,
			)
			if err != nil {
				tw.logger.Errorf("error getting candles: %v", err)
				continue // just ignoring it
			}

			formattedCandles := getFormattedCandles(candles)
			tw.logger.Debug("candles:", formattedCandles)

			if len(candles) < 6 {
				tw.logger.Warn("too few candles to proceed: expecting %d, got %d", 6, len(candles))
				continue // just skip
			}

			i := techan.NewTrendlineIndicator(techan.NewClosePriceIndicator(candlesToTimeSeries(candles)), len(candles)-3)
			trend := i.Calculate(len(candles) - 4)

			if trend.Float() < tw.config.TrendToTrade {
				tw.logger.Debugf("calculated trend lower when expected: %f < %f", trend.Float(), tw.config.TrendToTrade)
				continue // just wait until next turn
			}

			orderBook, err := services.MarketDataService.GetOrderBook(sdk.Figi(tw.Figi), 1)
			if err != nil {
				tw.logger.Errorf("error getting order Book: %v", err)
				continue // just ignoring it
			}

			tw.logger.Infof("last price: %d.%d", orderBook.LastPrice.Units, orderBook.LastPrice.Nano)
		case <-ctx.Done():
			tw.logger.Info("worker stopped!")
			return nil
		}
	}
}
