package gamble

import (
	"context"
	"errors"
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
	tradeutil "github.com/elkopass/BITA/internal/trade/util"
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

	sellFlag   bool           // if true, worker is trying to sell assets
	orderPrice *pb.MoneyValue // if order is set

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
					tw.orderPrice = nil
					tw.sellFlag = !tw.sellFlag
					go tw.checkPortfolio()
				} else {
					tw.logger.With("order_id", tw.orderID).Debug("order is still placed")
				}
				continue
			}

			if tw.sellFlag {
				tw.tryToSellInstrument()
			} else {
				tw.tryToBuyInstrument()
			}
		case <-ctx.Done():
			// TODO: implement sell logic on interrupt
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

	tw.logger.Info("positions: ", tradeutil.GetFormattedPositions(portfolio.Positions))

	for _, p := range portfolio.Positions {
		if p.CurrentPrice != nil {
			metrics.PortfolioPositionCurrentPrice.WithLabelValues(tw.accountID, tw.Figi).Set(tradeutil.MoneyValueToFloat(*p.CurrentPrice))
		}
		if p.ExpectedYield != nil {
			metrics.PortfolioPositionExpectedYield.WithLabelValues(tw.accountID, tw.Figi).Set(tradeutil.QuotationToFloat(*p.ExpectedYield))
		}
	}

	metrics.PortfolioInstrumentsAmount.WithLabelValues(tw.accountID, "bonds",
		portfolio.TotalAmountBonds.Currency).Set(tradeutil.MoneyValueToFloat(*portfolio.TotalAmountBonds))
	metrics.PortfolioInstrumentsAmount.WithLabelValues(tw.accountID, "currencies",
		portfolio.TotalAmountCurrencies.Currency).Set(tradeutil.MoneyValueToFloat(*portfolio.TotalAmountCurrencies))
	metrics.PortfolioInstrumentsAmount.WithLabelValues(tw.accountID, "etfs",
		portfolio.TotalAmountEtf.Currency).Set(tradeutil.MoneyValueToFloat(*portfolio.TotalAmountEtf))
	metrics.PortfolioInstrumentsAmount.WithLabelValues(tw.accountID, "futures",
		portfolio.TotalAmountFutures.Currency).Set(tradeutil.MoneyValueToFloat(*portfolio.TotalAmountFutures))
	metrics.PortfolioInstrumentsAmount.WithLabelValues(tw.accountID, "shares",
		portfolio.TotalAmountShares.Currency).Set(tradeutil.MoneyValueToFloat(*portfolio.TotalAmountShares))

	if portfolio.ExpectedYield != nil {
		tw.logger.Infof("expected yield: %d.%d", portfolio.ExpectedYield.Units, portfolio.ExpectedYield.Nano)
		metrics.PortfolioExpectedYieldOverall.WithLabelValues(tw.accountID).Set(tradeutil.QuotationToFloat(*portfolio.ExpectedYield))
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
			state.LotsExecuted, state.LotsRequested,
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
	if state.ExecutionReportStatus == pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_CANCELLED {
		tw.handleCancellation()
		return false
	}
	if state.ExecutionReportStatus == pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_REJECTED {
		tw.handleCancellation()
		return false
	}

	// all another cases are OK to place a new order
	direction := pb.OrderDirection_ORDER_DIRECTION_BUY.String()
	if tw.sellFlag {
		direction = pb.OrderDirection_ORDER_DIRECTION_SELL.String()
	}

	metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), tw.Figi, direction).Dec()
	metrics.OrdersFulfilled.WithLabelValues(loggy.GetBotID(), tw.Figi, direction).Inc()

	if tw.sellFlag {
		metrics.InstrumentsPurchased.WithLabelValues(loggy.GetBotID(), tw.Figi).Dec()
	} else {
		metrics.InstrumentsPurchased.WithLabelValues(loggy.GetBotID(), tw.Figi).Inc()
	}

	return true
}

func (tw *TradeWorker) tryToSellInstrument() {
	orderBook, err := services.MarketDataService.GetOrderBook(tw.Figi, 10)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		return // just ignoring it
	}

	priceIsOK := tw.priceIsOkToSell(*orderBook)
	if !priceIsOK {
		tw.logger.Debug("price is not OK to sell")
		return // wait for the next turn
	}

	order, err := services.SandboxService.PostSandboxOrder(
		&pb.PostOrderRequest{
			Figi:      tw.Figi,
			OrderId:   uuid.New().String(),
			Quantity:  int64(tw.config.LotsToBuy),
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
	tw.orderPrice = order.InitialOrderPrice
	tw.logger.With("order_id", tw.orderID).
		Infof("sell order created, initial price: %d.%d %s, current status: %s",
			order.InitialOrderPrice.Units, order.InitialOrderPrice.Nano,
			order.InitialOrderPrice.Currency, order.ExecutionReportStatus.String())

	metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), tw.Figi,
		pb.OrderDirection_ORDER_DIRECTION_SELL.String()).Inc()

	go tw.checkPortfolio()
}

func (tw *TradeWorker) tryToBuyInstrument() {
	trendIsOK, _ := tw.trendIsOkToBuy()
	if !trendIsOK {
		return // wait for the next turn
	}

	orderBook, err := services.MarketDataService.GetOrderBook(tw.Figi, 10)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		return // just ignoring it
	}

	closePrice := tradeutil.QuotationToFloat(*orderBook.ClosePrice)
	lastPrice := tradeutil.QuotationToFloat(*orderBook.LastPrice)
	limitUp := tradeutil.QuotationToFloat(*orderBook.LimitUp)
	limitDown := tradeutil.QuotationToFloat(*orderBook.LimitDown)

	metrics.InstrumentLastPrice.WithLabelValues(tw.Figi).Set(lastPrice)
	tw.logger.Infof("last price: %f, close price: %f limit up: %f, limit down: %f",
		lastPrice, closePrice, limitUp, limitDown)

	order, err := services.SandboxService.PostSandboxOrder(
		&pb.PostOrderRequest{
			Figi:      tw.Figi,
			OrderId:   uuid.New().String(),
			Quantity:  int64(tw.config.LotsToBuy),
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
	tw.orderPrice = order.InitialOrderPrice
	tw.logger.With("order_id", tw.orderID).
		Infof("buy order created, initial price: %d.%d %s, current status: %s",
			order.InitialOrderPrice.Units, order.InitialOrderPrice.Nano,
			order.InitialOrderPrice.Currency, order.ExecutionReportStatus.String())

	metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), tw.Figi,
		pb.OrderDirection_ORDER_DIRECTION_BUY.String()).Inc()
}

func (tw TradeWorker) tradingStatusIsOkToTrade() bool {
	status, err := services.MarketDataService.GetTradingStatus(tw.Figi)
	if err != nil {
		tw.logger.Errorf("error getting trading status: %v", err)
		return false
	}

	tw.logger.Infof("trading status: %s", status.TradingStatus.String())

	for _, s := range pb.SecurityTradingStatus_name {
		metrics.InstrumentTradingStatus.WithLabelValues(tw.Figi, s).Set(0)
	}
	metrics.InstrumentTradingStatus.WithLabelValues(tw.Figi, status.TradingStatus.String()).Set(1)

	return status.TradingStatus == pb.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING
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

	if len(shortCandles) < 6 || len(longCandles) < 6 {
		tw.logger.Warnf("too few candles to proceed: expecting at least %d, got %d and %d",
			6, len(shortCandles), len(longCandles))
		return false, nil
	}

	si := techan.NewTrendlineIndicator(techan.NewClosePriceIndicator(tradeutil.CandlesToTimeSeries(shortCandles)), len(shortCandles)-3)
	shortTrend := si.Calculate(len(shortCandles) - 4).Float()

	li := techan.NewTrendlineIndicator(techan.NewClosePriceIndicator(tradeutil.CandlesToTimeSeries(longCandles)), len(longCandles)-3)
	longTrend := li.Calculate(len(longCandles) - 4).Float()

	tw.logger.Debugf("calculated short trend: %f, expected: %f", shortTrend, tw.config.ShortTrendToTrade)
	tw.logger.Debugf("calculated long trend: %f, expected: %f", longTrend, tw.config.LongTrendToTrade)

	if longTrend > tw.config.LongTrendToTrade && shortTrend > tw.config.ShortTrendToTrade {
		return true, nil
	}

	return false, nil
}

func (tw *TradeWorker) priceIsOkToSell(orderBook pb.GetOrderBookResponse) bool {
	closePrice := tradeutil.QuotationToFloat(*orderBook.ClosePrice)
	lastPrice := tradeutil.QuotationToFloat(*orderBook.LastPrice)
	limitUp := tradeutil.QuotationToFloat(*orderBook.LimitUp)
	limitDown := tradeutil.QuotationToFloat(*orderBook.LimitDown)

	expectedProfit := closePrice * tw.config.TakeProfitCoef
	expectedLoss := closePrice * tw.config.StopLossCoef

	metrics.InstrumentLastPrice.WithLabelValues(tw.Figi).Set(lastPrice)
	tw.logger.Infof("expected price: %f, last: %f, close: %f limit up: %f, limit down: %f, stop loss: %f",
		expectedProfit, lastPrice, closePrice, limitUp, limitDown, expectedLoss)

	if closePrice < expectedLoss || closePrice > expectedProfit {
		return true
	}
	return false
}

func (tw *TradeWorker) handleCancellation() {
	tw.logger.With("order_id", tw.orderID).Warn("order is cancelled")

	tw.orderID = ""
	tw.orderPrice = nil
}
