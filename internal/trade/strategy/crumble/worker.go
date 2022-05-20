package crumble

import (
	"context"
	"errors"
	"github.com/elkopass/BITA/internal/config"
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
	cb "github.com/elkopass/BITA/internal/trade/breaker"
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

	whichAverageIsBigger string

	sellFlag        bool           // if true, worker is trying to sell assets
	orderPrice      *pb.MoneyValue // if order is set
	orderPlacedTime *int64         // if order is set

	logger  *zap.SugaredLogger
	breaker cb.CircuitBreaker
	config  TradeConfig
}

func NewTradeWorker(figi, accountID string) *TradeWorker {
	id := strings.Split(uuid.New().String(), "-")[0]

	return &TradeWorker{
		ID:        id,
		Figi:      figi,
		accountID: accountID,
		config:    *NewTradeConfig(),
		breaker:   *cb.NewCircuitBreaker(),
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
			if tw.breaker.WorkerMustExit() {
				tw.logger.Error("worker stopped by circuit breaker")
				metrics.StoppedByCircuitBreaker.WithLabelValues(loggy.GetBotID(), tw.Figi).Inc()
				return
			}

			if !tw.tradingStatusIsOkToTrade() {
				continue // just skip
			}

			if tw.orderID != "" {
				if tw.orderIsFulfilled() {
					tw.orderID = ""
					tw.sellFlag = !tw.sellFlag
					go tw.checkPortfolio()
				} else {
					tw.logger.With("order_id", tw.orderID).Debug("order is still placed")
					tw.checkNeedForCancel()
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

			if config.TradeBotConfig().SellOnExit && tw.sellFlag {
				tw.logger.Info("SELL_ON_EXIT flag is set, trying to sell an asset...")
				return tw.sellOnExit()
			}

			return nil
		}
	}
}

// sellOnExit immediately creates sell order if worker has an instrument.
func (tw TradeWorker) sellOnExit() error {
	orderBook, err := services.MarketDataService.GetOrderBook(tw.Figi, 10)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		tw.breaker.IncFailures()
		return err
	}

	fairPrice, err := tradeutil.CalculateFairBuyPrice(*orderBook)
	if err != nil {
		tw.logger.Errorf("can not calculate fair price: %v", err)
		return err
	}

	orderRequest := &pb.PostOrderRequest{
		Figi:      tw.Figi,
		OrderId:   uuid.New().String(),
		Quantity:  int64(tw.config.LotsToBuy),
		Price:     fairPrice,
		AccountId: tw.accountID,
		OrderType: pb.OrderType_ORDER_TYPE_MARKET,
		Direction: pb.OrderDirection_ORDER_DIRECTION_SELL,
	}

	var orderResponse *pb.PostOrderResponse
	if config.TradeBotConfig().IsSandbox {
		orderResponse, err = services.SandboxService.PostSandboxOrder(orderRequest)
	} else {
		orderResponse, err = services.OrdersService.PostOrder(orderRequest)
	}

	if err != nil {
		tw.logger.Errorf("can not post sell order: %v", err)
		tw.breaker.IncFailures()
		return err
	}

	tw.orderID = orderResponse.OrderId
	tw.orderPrice = &pb.MoneyValue{
		Units:    fairPrice.Units,
		Nano:     fairPrice.Nano,
		Currency: orderResponse.InitialOrderPrice.Currency,
	}

	tw.logger.With("order_id", tw.orderID).
		Infof("sell order created, fair price: %d.%d, initial price: %d.%d %s, current status: %s",
			fairPrice.Units, fairPrice.Nano,
			orderResponse.InitialOrderPrice.Units, orderResponse.InitialOrderPrice.Nano,
			orderResponse.InitialOrderPrice.Currency, orderResponse.ExecutionReportStatus.String())

	metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), tw.Figi,
		pb.OrderDirection_ORDER_DIRECTION_SELL.String()).Inc()

	return nil
}

// checkPortfolio calls sdk.OperationsService.GetPortfolio and updates portfolio metrics.
func (tw *TradeWorker) checkPortfolio() {
	var portfolio *pb.PortfolioResponse
	var err error

	if config.TradeBotConfig().IsSandbox {
		portfolio, err = services.SandboxService.GetSandboxPortfolio(tw.accountID)
	} else {
		portfolio, err = services.OperationsService.GetPortfolio(tw.accountID)
	}

	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		tw.breaker.IncFailures()
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

// orderIsFulfilled calls sdk.OrdersService.GetOrderState and checks ExecutionReportStatus.
// If order is not fulfilled, it will return false or even call the handleCancellation.
func (tw *TradeWorker) orderIsFulfilled() bool {
	var state *pb.OrderState
	var err error

	if config.TradeBotConfig().IsSandbox {
		state, err = services.SandboxService.GetSandboxOrderState(tw.accountID, tw.orderID)
	} else {
		state, err = services.OrdersService.GetOrderState(tw.accountID, tw.orderID)
	}

	if err != nil {
		tw.logger.With("order_id", tw.orderID).Errorf("can not check order state: %v", err)
		tw.breaker.IncFailures()
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

// tryToSellInstrument calls sdk.MarketDataService.GetOrderBook and if priceIsOkToSell
// the order will be placed and orderID will be set along with orderPrice.
func (tw *TradeWorker) tryToSellInstrument() {
	orderBook, err := services.MarketDataService.GetOrderBook(tw.Figi, 10)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		tw.breaker.IncFailures()
		return // just ignoring it
	}

	priceIsOK := tw.priceIsOkToSell(*orderBook)
	indicatorIsOK, _ := tw.indicatorIsOkToSell()
	if !priceIsOK && !indicatorIsOK {
		return // wait for the next turn
	}

	fairPrice, err := tradeutil.CalculateFairSellPrice(*orderBook)
	if err != nil {
		tw.logger.Errorf("can not calculate fair price: %v", err)
		return // try again next time
	}

	orderRequest := &pb.PostOrderRequest{
		Figi:      tw.Figi,
		OrderId:   uuid.New().String(),
		Quantity:  int64(tw.config.LotsToBuy),
		Price:     fairPrice,
		AccountId: tw.accountID,
		OrderType: pb.OrderType_ORDER_TYPE_LIMIT,
		Direction: pb.OrderDirection_ORDER_DIRECTION_SELL,
	}

	var orderResponse *pb.PostOrderResponse
	if config.TradeBotConfig().IsSandbox {
		orderResponse, err = services.SandboxService.PostSandboxOrder(orderRequest)
	} else {
		orderResponse, err = services.OrdersService.PostOrder(orderRequest)
	}

	if err != nil {
		tw.logger.Errorf("can not post sell order: %v", err)
		tw.breaker.IncFailures()
		return // nothing bad happened, let's proceed
	}

	tw.orderID = orderResponse.OrderId
	tw.orderPrice = &pb.MoneyValue{
		Units:    fairPrice.Units,
		Nano:     fairPrice.Nano,
		Currency: orderResponse.InitialOrderPrice.Currency,
	}

	t := time.Now().Unix()
	tw.orderPlacedTime = &t

	tw.logger.With("order_id", tw.orderID).
		Infof("sell order created, fair price: %d.%d, initial price: %d.%d %s, current status: %s",
			fairPrice.Units, fairPrice.Nano,
			orderResponse.InitialOrderPrice.Units, orderResponse.InitialOrderPrice.Nano,
			orderResponse.InitialOrderPrice.Currency, orderResponse.ExecutionReportStatus.String())

	metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), tw.Figi,
		pb.OrderDirection_ORDER_DIRECTION_SELL.String()).Inc()

	go tw.checkPortfolio()
}

// tryToSellInstrument calls sdk.MarketDataService.GetOrderBook and if indicatorIsOkToBuy
// the order will be placed and orderID will be set along with orderPrice.
func (tw *TradeWorker) tryToBuyInstrument() {
	indicatorIsOK, _ := tw.indicatorIsOkToBuy()
	if !indicatorIsOK {
		return // wait for the next turn
	}

	orderBook, err := services.MarketDataService.GetOrderBook(tw.Figi, 10)
	if err != nil {
		tw.logger.Errorf("error getting order book: %v", err)
		tw.breaker.IncFailures()
		return // just ignoring it
	}

	fairPrice, err := tradeutil.CalculateFairBuyPrice(*orderBook)
	if err != nil {
		tw.logger.Errorf("can not calculate fair price: %v", err)
		return // try again next time
	}

	closePrice := tradeutil.QuotationToFloat(*orderBook.ClosePrice)
	lastPrice := tradeutil.QuotationToFloat(*orderBook.LastPrice)
	fairMarketPrice := tradeutil.QuotationToFloat(*fairPrice)

	metrics.InstrumentLastPrice.WithLabelValues(tw.Figi).Set(lastPrice)
	metrics.InstrumentFairPrice.WithLabelValues(tw.Figi).Set(fairMarketPrice)
	tw.logger.Infof("last price: %f, close price: %f, fair price: %f",
		lastPrice, closePrice, fairMarketPrice)

	orderRequest := &pb.PostOrderRequest{
		Figi:      tw.Figi,
		OrderId:   uuid.New().String(),
		Quantity:  int64(tw.config.LotsToBuy),
		Price:     fairPrice,
		AccountId: tw.accountID,
		OrderType: pb.OrderType_ORDER_TYPE_LIMIT,
		Direction: pb.OrderDirection_ORDER_DIRECTION_BUY,
	}

	var orderResponse *pb.PostOrderResponse
	if config.TradeBotConfig().IsSandbox {
		orderResponse, err = services.SandboxService.PostSandboxOrder(orderRequest)
	} else {
		orderResponse, err = services.OrdersService.PostOrder(orderRequest)
	}
	if err != nil {
		tw.logger.Errorf("can not post buy order: %v", err)
		return // nothing bad happened, let's proceed
	}

	tw.orderID = orderResponse.OrderId
	tw.orderPrice = &pb.MoneyValue{
		Units:    fairPrice.Units,
		Nano:     fairPrice.Nano,
		Currency: orderResponse.InitialOrderPrice.Currency,
	}

	t := time.Now().Unix()
	tw.orderPlacedTime = &t

	tw.logger.With("order_id", tw.orderID).
		Infof("buy order created, fair price: %d.%d, initial price: %d.%d %s, current status: %s",
			fairPrice.Units, fairPrice.Nano,
			orderResponse.InitialOrderPrice.Units, orderResponse.InitialOrderPrice.Nano,
			orderResponse.InitialOrderPrice.Currency, orderResponse.ExecutionReportStatus.String())

	metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), tw.Figi,
		pb.OrderDirection_ORDER_DIRECTION_BUY.String()).Inc()
}

// tradingStatusIsOkToTrade returns true if trading status is normal.
func (tw TradeWorker) tradingStatusIsOkToTrade() bool {
	status, err := services.MarketDataService.GetTradingStatus(tw.Figi)
	if err != nil {
		tw.logger.Errorf("error getting trading status: %v", err)
		tw.breaker.IncFailures()
		return false
	}

	tw.logger.Infof("trading status: %s", status.TradingStatus.String())
	for _, s := range pb.SecurityTradingStatus_name {
		metrics.InstrumentTradingStatus.WithLabelValues(tw.Figi, s).Set(0)
	}
	metrics.InstrumentTradingStatus.WithLabelValues(tw.Figi, status.TradingStatus.String()).Set(1)

	return status.TradingStatus == pb.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING
}

// indicatorIsOkToBuy checks MA-indicator and returns true if it's OK to buy.
func (tw *TradeWorker) indicatorIsOkToBuy() (bool, error) {
	candles, err := services.MarketDataService.GetCandles(
		tw.Figi,
		timestamppb.New(time.Now().Add(-time.Duration(tw.config.CandlesIntervalHours)*time.Hour)),
		timestamppb.Now(),
		pb.CandleInterval_CANDLE_INTERVAL_HOUR,
	)

	if err != nil {
		tw.breaker.IncFailures()
		return false, errors.New("error getting short candles: " + err.Error())
	}

	minCandles := tw.config.LongWindow + 1
	if len(candles) < minCandles {
		tw.logger.Warnf("too few candles to proceed: expecting at least %d, got %d and %d",
			minCandles, len(candles))
		return false, nil
	}

	si := techan.NewMMAIndicator(techan.NewClosePriceIndicator(tradeutil.CandlesToTimeSeries(candles)), tw.config.ShortWindow)
	li := techan.NewMMAIndicator(techan.NewClosePriceIndicator(tradeutil.CandlesToTimeSeries(candles)), tw.config.LongWindow)

	shortMMA := si.Calculate(0).Float()
	longMMA := li.Calculate(0).Float()

	tw.logger.Infof("calculated short MMA: %f", shortMMA)
	tw.logger.Infof("calculated long MMA: %f", longMMA)

	if shortMMA > longMMA {
		tw.whichAverageIsBigger = "short"
		return false, nil
	}
	if shortMMA < longMMA {
		tw.whichAverageIsBigger = "long"
		return false, nil
	}

	if shortMMA == longMMA && tw.whichAverageIsBigger == "short" {
		tw.whichAverageIsBigger = "long"
		return true, nil
	}

	return false, nil
}

// indicatorIsOkToSell checks MA-indicator once again and returns true if it's OK to sell.
func (tw *TradeWorker) indicatorIsOkToSell() (bool, error) {
	candles, err := services.MarketDataService.GetCandles(
		tw.Figi,
		timestamppb.New(time.Now().Add(-time.Duration(tw.config.CandlesIntervalHours)*time.Hour)),
		timestamppb.Now(),
		pb.CandleInterval_CANDLE_INTERVAL_HOUR,
	)
	if err != nil {
		tw.breaker.IncFailures()
		return false, errors.New("error getting short candles: " + err.Error())
	}

	minCandles := tw.config.LongWindow + 1
	if len(candles) < minCandles {
		tw.logger.Warnf("too few candles to proceed: expecting at least %d, got %d and %d",
			minCandles, len(candles))
		return false, nil
	}

	si := techan.NewMMAIndicator(techan.NewClosePriceIndicator(tradeutil.CandlesToTimeSeries(candles)), tw.config.ShortWindow)
	li := techan.NewMMAIndicator(techan.NewClosePriceIndicator(tradeutil.CandlesToTimeSeries(candles)), tw.config.LongWindow)

	shortMMA := si.Calculate(0).Float()
	longMMA := li.Calculate(0).Float()

	tw.logger.Debugf("calculated short MMA: %f", shortMMA)
	tw.logger.Debugf("calculated long MMA: %f", longMMA)

	if shortMMA > longMMA {
		tw.whichAverageIsBigger = "short"
		return false, nil
	}
	if shortMMA < longMMA {
		tw.whichAverageIsBigger = "long"
		return false, nil
	}

	if shortMMA == longMMA && tw.whichAverageIsBigger == "long" {
		tw.whichAverageIsBigger = "short"
		return true, nil
	}

	return false, nil
}

// priceIsOkToSell returns true if (price > expected profit) or (price < expected loss).
func (tw *TradeWorker) priceIsOkToSell(orderBook pb.GetOrderBookResponse) bool {
	fairPrice, err := tradeutil.CalculateFairSellPrice(orderBook)
	if err != nil {
		tw.logger.Warnf("can't calculate fairMarketPrice price: %v", err.Error())
		return false
	}

	lastOrderPrice := tradeutil.MoneyValueToFloat(*tw.orderPrice)
	closePrice := tradeutil.QuotationToFloat(*orderBook.ClosePrice)
	lastPrice := tradeutil.QuotationToFloat(*orderBook.LastPrice)
	fairMarketPrice := tradeutil.QuotationToFloat(*fairPrice)

	expectedProfit := lastOrderPrice * tw.config.TakeProfitCoef
	expectedLoss := lastOrderPrice * tw.config.StopLossCoef

	metrics.InstrumentLastPrice.WithLabelValues(tw.Figi).Set(lastPrice)
	metrics.InstrumentFairPrice.WithLabelValues(tw.Figi).Set(fairMarketPrice)
	tw.logger.Infof("order price: %f, fair price: %f, expected: %f, last: %f, close: %f, stop loss: %f",
		lastOrderPrice, fairMarketPrice, expectedProfit, lastPrice, closePrice, expectedLoss)

	if fairMarketPrice < expectedLoss {
		metrics.StopLossDecisions.WithLabelValues(loggy.GetBotID(), tw.Figi).Inc()
		return true
	}
	if fairMarketPrice > expectedProfit {
		metrics.TakeProfitDecisions.WithLabelValues(loggy.GetBotID(), tw.Figi).Inc()
		return true
	}

	return false
}

// checkNeedForCancel tries to cancel orders older than TradeConfig.SecondsToCancelOrder.
func (tw *TradeWorker) checkNeedForCancel() {
	if *tw.orderPlacedTime-time.Now().Unix() > tw.config.SecondsToCancelOrder {
		_, err := services.OrdersService.CancelOrder(tw.accountID, tw.orderID)
		if err != nil {
			tw.logger.Warnf("can not cancel order: %v", err)
			return
		}

		tw.handleCancellation()
	}
}

// handleCancellation unsets orderID.
func (tw *TradeWorker) handleCancellation() {
	metrics.OrdersCancelled.WithLabelValues(loggy.GetBotID(), tw.Figi).Inc()

	tw.logger.With("order_id", tw.orderID).Warn("order is cancelled")
	tw.orderPlacedTime = nil
	tw.orderID = ""
}
