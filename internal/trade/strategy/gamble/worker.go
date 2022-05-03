package gamble

import (
	"context"
	"errors"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/google/uuid"
	"github.com/sdcoffey/techan"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"sync"
	"time"
)

type TradeWorker struct {
	ID        string
	Figi      string
	orderID   string
	accountID string

	sellFlag bool // if true, bot is trying to sell his assets

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
		sellFlag:  false,
		logger: loggy.GetLogger().Sugar().
			With("bot_id", loggy.GetBotID()).
			With("account_id", accountID).
			With("worker_id", id).
			With("figi", figi),
	}
}

func (tw TradeWorker) Run(ctx context.Context, wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	tw.logger = tw.logger.With("sell_flag", tw.sellFlag)
	tw.logger.Debug("start trading...")

	for {
		select {
		case <-time.After(time.Duration(tw.config.WorkerSleepDurationSeconds) * time.Second):
			if !tw.tradingStatusIsOkToTrade() {
				continue // just skip
			}

			if tw.orderID != "" {
				if tw.orderIsFulfilled() {
					tw.orderID = ""
					tw.sellFlag = !tw.sellFlag
					tw.checkPortfolio()
				} else {
					tw.logger.With("order_id", tw.orderID).Debug("order is still placed")
					tw.checkInstrument()
				}
				continue
			}

			if tw.sellFlag {
				tw.tryToSellInstrument()
			} else {
				tw.tryToBuyInstrument()
			}
		case <-ctx.Done():
			tw.logger.Info("worker stopped!")
			return nil
		}
	}
}

func (tw *TradeWorker) checkPortfolio() {
	portfolio, err := services.SandboxService.GetSandboxPortfolio(tw.accountID)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		return // just ignoring it
	}

	tw.logger.Info("positions: ", getFormattedPositions(portfolio.Positions))
	if portfolio.ExpectedYield != nil {
		tw.logger.Infof("expected yield: %d.%d", portfolio.ExpectedYield.Units, portfolio.ExpectedYield.Nano)
	}
}

func (tw *TradeWorker) orderIsFulfilled() bool {
	state, err := services.SandboxService.GetSandboxOrderState(tw.accountID, tw.orderID)
	if err != nil {
		tw.logger.With("order_id", tw.orderID).Errorf("can not check order state: %v", err)
		return false
	}

	tw.logger.With("order_id", tw.orderID).
		Infof("order status: %s, fulfilled %d/%d, current price: %d.%d %s",
			state.ExecutionReportStatus.String(),
			state.LotsExecuted,
			state.LotsRequested,
			state.AveragePositionPrice.Units,
			state.AveragePositionPrice.Nano,
			state.AveragePositionPrice.Currency,
		)

	tw.logger.With("order_id", tw.orderID).
		Infof("execution status: %s", state.ExecutionReportStatus)

	if state.ExecutionReportStatus == pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_NEW {
		return false
	}
	if state.ExecutionReportStatus == pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_PARTIALLYFILL {
		return false
	}

	// all another cases are OK to place a new order
	return true
}

func (tw *TradeWorker) checkInstrument() {
	orderBook, err := services.MarketDataService.GetOrderBook(tw.Figi, 1)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		return // just ignoring it
	}

	tw.logger.Infof("last price: %d.%d, close price: %d.%d, limit up: %d.%d, limit down: %d.%d",
		orderBook.LastPrice.Units, orderBook.LastPrice.Nano,
		orderBook.ClosePrice.Units, orderBook.ClosePrice.Nano,
		orderBook.LimitUp.Units, orderBook.LimitUp.Nano,
		orderBook.LimitDown.Units, orderBook.LimitDown.Nano,
	)
}

func (tw *TradeWorker) tryToSellInstrument() {
	priceIsOK, err := tw.priceIsOkToSell()
	if err != nil {
		tw.logger.Errorf("error getting price: %v", err)
		return // just ignore
	}
	if !priceIsOK {
		tw.logger.Debug("price is not OK to sell")
		return // wait for the next turn
	}

	orderBook, err := services.MarketDataService.GetOrderBook(tw.Figi, 10)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		return // just ignoring it
	}
	tw.logger.Infof("last price: %d.%d, close price: %d.%d, limit up: %d.%d, limit down: %d.%d",
		orderBook.LastPrice.Units, orderBook.LastPrice.Nano,
		orderBook.ClosePrice.Units, orderBook.ClosePrice.Nano,
		orderBook.LimitUp.Units, orderBook.LimitUp.Nano,
		orderBook.LimitDown.Units, orderBook.LimitDown.Nano,
	)

	order, err := services.SandboxService.PostSandboxOrder(
		&pb.PostOrderRequest{
			Figi:      tw.Figi,
			OrderId:   uuid.New().String(),
			Quantity:  int64(tw.config.AmountToBuy),
			Price:     orderBook.LastPrice,
			AccountId: tw.accountID,
			OrderType: pb.OrderType_ORDER_TYPE_LIMIT,
			Direction: pb.OrderDirection_ORDER_DIRECTION_SELL,
		},
	)
	if err != nil {
		tw.logger.Errorf("can not post sell order: %v", err)
		return // nothing bad happened, let's proceed
	}

	tw.orderID = order.OrderId
	tw.logger.With("order_id", tw.orderID).
		Infof("sell order created, current status: %s", order.ExecutionReportStatus.String())

	tw.checkPortfolio()
}

func (tw *TradeWorker) tryToBuyInstrument() {
	trendIsOK, err := tw.trendIsOkToBuy()
	if !trendIsOK {
		return // wait for the next turn
	}

	orderBook, err := services.MarketDataService.GetOrderBook(tw.Figi, 10)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		return // just ignoring it
	}
	tw.logger.Infof("last price: %d.%d, close price: %d.%d, limit up: %d.%d, limit down: %d.%d",
		orderBook.LastPrice.Units, orderBook.LastPrice.Nano,
		orderBook.ClosePrice.Units, orderBook.ClosePrice.Nano,
		orderBook.LimitUp.Units, orderBook.LimitUp.Nano,
		orderBook.LimitDown.Units, orderBook.LimitDown.Nano,
	)

	order, err := services.SandboxService.PostSandboxOrder(
		&pb.PostOrderRequest{
			Figi:      tw.Figi,
			OrderId:   uuid.New().String(),
			Quantity:  int64(tw.config.AmountToBuy),
			Price:     orderBook.LastPrice,
			AccountId: tw.accountID,
			OrderType: pb.OrderType_ORDER_TYPE_LIMIT,
			Direction: pb.OrderDirection_ORDER_DIRECTION_BUY,
		},
	)
	if err != nil {
		tw.logger.Errorf("can not post buy order: %v", err)
		return // nothing bad happened, let's proceed
	}

	tw.orderID = order.OrderId
	tw.logger.With("order_id", tw.orderID).
		Infof("buy order created, current status: %s", order.ExecutionReportStatus.String())
}

func (tw TradeWorker) tradingStatusIsOkToTrade() bool {
	status, err := services.MarketDataService.GetTradingStatus(tw.Figi)
	if err != nil {
		tw.logger.Errorf("error getting trading status: %v", err)
		return false
	}

	tw.logger.Infof("trading status: %s", status.TradingStatus.String())
	if status.TradingStatus == pb.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING {
		return true
	}

	return false
}

func (tw *TradeWorker) trendIsOkToBuy() (bool, error) {
	shortCandles, err := services.MarketDataService.GetCandles(
		tw.Figi,
		timestamppb.New(time.Now().Add(-time.Duration(tw.config.ShortTrendIntervalSeconds)*time.Second)),
		timestamppb.Now(),
		pb.CandleInterval_CANDLE_INTERVAL_1_MIN,
	)
	if err != nil {
		return false, errors.New("error getting short candles: " + err.Error())
	}

	longCandles, err := services.MarketDataService.GetCandles(
		tw.Figi,
		timestamppb.New(time.Now().Add(-time.Duration(tw.config.LongTrendIntervalSeconds)*time.Second)),
		timestamppb.Now(),
		pb.CandleInterval_CANDLE_INTERVAL_5_MIN,
	)
	if err != nil {
		return false, errors.New("error getting long candles: " + err.Error())
	}

	//formattedShortCandles := getFormattedCandles(shortCandles)
	//tw.logger.Debug("short candles:", formattedShortCandles)

	//formattedLongCandles := getFormattedCandles(longCandles)
	//tw.logger.Debug("short candles:", formattedLongCandles)

	if len(shortCandles) < 6 || len(longCandles) < 6 {
		tw.logger.Warnf("too few candles to proceed: expecting at least %d, got %d and %d",
			6, len(shortCandles), len(longCandles))
		return false, nil
	}

	si := techan.NewTrendlineIndicator(techan.NewClosePriceIndicator(candlesToTimeSeries(shortCandles)), len(shortCandles)-3)
	shortTrend := si.Calculate(len(shortCandles) - 4).Float()

	li := techan.NewTrendlineIndicator(techan.NewClosePriceIndicator(candlesToTimeSeries(longCandles)), len(longCandles)-3)
	longTrend := li.Calculate(len(longCandles) - 4).Float()

	tw.logger.Debugf("calculated short trend: %f, expected: %f", shortTrend, tw.config.ShortTrendToTrade)
	tw.logger.Debugf("calculated long trend: %f, expected: %f", longTrend, tw.config.LongTrendToTrade)

	if longTrend > tw.config.LongTrendToTrade && shortTrend > tw.config.ShortTrendToTrade {
		return true, nil
	}

	return false, nil
}

func (tw *TradeWorker) priceIsOkToSell() (bool, error) {
	// TODO: implementation
	return true, nil
}
