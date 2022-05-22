package tumble

import (
	"context"
	"errors"
	"fmt"
	"github.com/elkopass/BITA/internal/config"
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/elkopass/BITA/internal/sdk"
	"github.com/elkopass/BITA/internal/trade/common"
	tradeutil "github.com/elkopass/BITA/internal/trade/util"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strings"
	"time"
)

type TradeBot struct {
	accountID string
	orders    map[string]Order // figi == key
	config    TradeConfig
	logger    *zap.SugaredLogger

	tradesStream sdk.OrdersStream
}

type Order struct {
	SellFlag        bool           // if true, worker is trying to sell assets
	OrderID         string         // if order is set
	OrderPrice      *pb.MoneyValue // if order is set
	OrderPlacedTime *int64         // if order is set
}

func NewTradeBot() *TradeBot {
	return &TradeBot{
		orders: make(map[string]Order),
		config: *NewTradeConfig(),
		logger: loggy.GetLogger().Sugar().With("bot_id", loggy.GetBotID()),
	}
}

func (tb TradeBot) Run(ctx context.Context) (err error) {
	tb.logger.Infof("starting with %s strategy and sdk v%s", config.TradeBotConfig().Strategy, sdk.Version)

	err = tb.setAccountID()
	if err != nil {
		return err
	}

	if config.TradeBotConfig().IsSandbox {
		return errors.New("strategy is not available in sandbox, " +
			"see https://github.com/Tinkoff/investAPI/issues/176")
	} else {
		tb.tradesStream = *sdk.NewOrdersStream(&pb.TradesStreamRequest{Accounts: []string{tb.accountID}})
	}

	figi := config.TradeBotConfig().Figi

	var instruments []*pb.OrderBookInstrument
	for _, f := range figi {
		instruments = append(instruments, &pb.OrderBookInstrument{Figi: f, Depth: int32(tb.config.OrderBookDepth)})
	}

	mds := sdk.NewMarketDataStream()

	request := pb.SubscribeOrderBookRequest{
		Instruments:        instruments,
		SubscriptionAction: pb.SubscriptionAction_SUBSCRIPTION_ACTION_SUBSCRIBE,
	}
	payload := &pb.MarketDataRequest_SubscribeOrderBookRequest{SubscribeOrderBookRequest: &request}

	err = mds.Send(&pb.MarketDataRequest{Payload: payload})
	if err != nil {
		return err
	}

	tradeStreamCtx, cancelTradeStreamListener := context.WithCancel(context.Background())
	go tb.listenTradeStream(tradeStreamCtx)

	for {
		msg, err := mds.Recv()
		if err != nil {
			tb.logger.Error(err)
		}

		orderBook := msg.GetOrderbook()
		if orderBook == nil {
			continue
		}

		tb.logger.Debug(tradeutil.GetFormattedOrderBook(orderBook))
		tb.makeDecision(orderBook)

		select {
		case <-time.After(1 * time.Millisecond):
			// pass
		case <-ctx.Done():
			cancelTradeStreamListener()

			// TODO: implement sell logic on interrupt
			if config.TradeBotConfig().IsSandbox {
				err = services.SandboxService.CloseSandboxAccount(tb.accountID)
				if err != nil {
					tb.logger.Errorf("can't close an account: %v", err)
				}
				tb.logger.Infof("account with ID %s closed successfully", tb.accountID)
			}

			tb.logger.Info("bot stopped!")
			return nil
		}
	}
}

// setAccountID gets account ID from config or creates a new one in sandbox.
func (tb *TradeBot) setAccountID() error {
	accountID := config.TradeBotConfig().AccountID
	if config.TradeBotConfig().IsSandbox {
		accountID, err := services.SandboxService.OpenSandboxAccount()
		if err != nil {
			return fmt.Errorf("can not create account: %v", err)
		}
		tb.logger.Infof("created new account with ID %s", accountID)

		tb.logger = tb.logger.With("account_id", accountID)
		tb.accountID = accountID
	} else {
		info, err := services.UsersService.GetInfo()
		if err != nil {
			return fmt.Errorf("can not get user info: %v", err)
		}
		tb.logger.Infof("user tariff: %s, qualified for work with %s",
			info.Tariff, strings.Join(info.QualifiedForWorkWith, ","))

		tb.logger = tb.logger.With("account_id", accountID)
		tb.accountID = accountID
	}

	return nil
}

// listenTradeStream receives fulfilled orders from stream.
func (tb *TradeBot) listenTradeStream(ctx context.Context) {
	for {
		msg, err := tb.tradesStream.Recv()
		if err != nil {
			tb.logger.Error(err)
		}

		orderTrades := msg.GetOrderTrades()
		if orderTrades != nil {
			tb.logger.With("order_id", orderTrades.OrderId).
				With("figi", orderTrades.Figi).
				Info("order is fulfilled")

			metrics.OrdersFulfilled.WithLabelValues(loggy.GetBotID(),
				orderTrades.Figi, orderTrades.Direction.String()).Inc()
			metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), orderTrades.Figi,
				orderTrades.Direction.String()).Dec()

			switch orderTrades.Direction {
			case pb.OrderDirection_ORDER_DIRECTION_BUY:
				metrics.InstrumentsPurchased.WithLabelValues(loggy.GetBotID(), orderTrades.Figi).Inc()
			case pb.OrderDirection_ORDER_DIRECTION_SELL:
				metrics.InstrumentsPurchased.WithLabelValues(loggy.GetBotID(), orderTrades.Figi).Dec()
			}

			delete(tb.orders, orderTrades.Figi)
			go tb.checkPortfolio()
		}

		select {
		case <-time.After(1000 * time.Millisecond):
			// pass
		case <-ctx.Done():
			tb.logger.Debug("stop trade stream listener")
			return
		}
	}
}

// makeDecision checks pb.OrderBook volumes with the goal to create buy/sell order.
func (tb *TradeBot) makeDecision(orderBook *pb.OrderBook) {
	var asksQuantity float64
	for _, ask := range orderBook.Asks {
		asksQuantity += float64(ask.Quantity)
	}

	var bidsQuantity float64
	for _, bid := range orderBook.Bids {
		bidsQuantity += float64(bid.Quantity)
	}

	tb.logger.Debugf("ask/bids ratio: %f, expected: %f",
		asksQuantity/bidsQuantity, tb.config.AsksBidsRatio)
	tb.logger.Debugf("bids/asks ratio: %f, expected: %f",
		bidsQuantity/asksQuantity, tb.config.BidsAsksRatio)

	if order, ok := tb.orders[orderBook.Figi]; ok {
		tb.logger.With("order_id", order.OrderID).Debug("order already exists")
		return
	}

	if bidsQuantity/asksQuantity > tb.config.BidsAsksRatio {
		tb.tryToBuy(orderBook)
	}
	if asksQuantity/bidsQuantity > tb.config.AsksBidsRatio {
		tb.tryToSell(orderBook)
	}
}

// tryToBuy tries to create buy order with price calculated on pb.OrderBook.
func (tb *TradeBot) tryToBuy(orderBook *pb.OrderBook) {
	fairPrice := orderBook.Bids[tb.config.OrderBookFairBidDepth].Price
	fairMarketPrice := tradeutil.QuotationToFloat(*fairPrice)

	metrics.InstrumentFairPrice.WithLabelValues(orderBook.Figi).Set(fairMarketPrice)
	tb.logger.Infof("fair price: %f", fairMarketPrice)

	orderRequest := &pb.PostOrderRequest{
		Figi:      orderBook.Figi,
		OrderId:   uuid.New().String(),
		Quantity:  int64(tb.config.LotsToBuy),
		Price:     fairPrice,
		AccountId: tb.accountID,
		OrderType: pb.OrderType_ORDER_TYPE_MARKET,
		Direction: pb.OrderDirection_ORDER_DIRECTION_BUY,
	}

	var orderResponse *pb.PostOrderResponse
	var err error

	if config.TradeBotConfig().IsSandbox {
		orderResponse, err = services.SandboxService.PostSandboxOrder(orderRequest)
	} else {
		orderResponse, err = services.OrdersService.PostOrder(orderRequest)
	}
	if err != nil {
		tb.logger.Errorf("can not post buy order: %v", err)
		return // nothing bad happened, let's proceed
	}

	var order Order

	order.SellFlag = false
	order.OrderID = orderResponse.OrderId
	order.OrderPrice = &pb.MoneyValue{
		Units:    fairPrice.Units,
		Nano:     fairPrice.Nano,
		Currency: orderResponse.InitialOrderPrice.Currency,
	}

	t := time.Now().Unix()
	order.OrderPlacedTime = &t

	tb.logger.With("order_id", order.OrderID).
		Infof("buy order created, fair price: %d.%d, initial price: %d.%d %s, current status: %s",
			fairPrice.Units, fairPrice.Nano,
			orderResponse.InitialOrderPrice.Units, orderResponse.InitialOrderPrice.Nano,
			orderResponse.InitialOrderPrice.Currency, orderResponse.ExecutionReportStatus.String())

	metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), orderBook.Figi,
		pb.OrderDirection_ORDER_DIRECTION_BUY.String()).Inc()

	tb.orders[orderBook.Figi] = order
}

// tryToSell tries to create sell order with price calculated on pb.OrderBook.
func (tb *TradeBot) tryToSell(orderBook *pb.OrderBook) {
	fairPrice := orderBook.Asks[5].Price
	fairMarketPrice := tradeutil.QuotationToFloat(*fairPrice)

	metrics.InstrumentFairPrice.WithLabelValues(orderBook.Figi).Set(fairMarketPrice)
	tb.logger.Infof("fair price: %f", fairMarketPrice)

	orderRequest := &pb.PostOrderRequest{
		Figi:      orderBook.Figi,
		OrderId:   uuid.New().String(),
		Quantity:  int64(tb.config.LotsToBuy),
		Price:     fairPrice,
		AccountId: tb.accountID,
		OrderType: pb.OrderType_ORDER_TYPE_MARKET,
		Direction: pb.OrderDirection_ORDER_DIRECTION_SELL,
	}

	var orderResponse *pb.PostOrderResponse
	var err error

	if config.TradeBotConfig().IsSandbox {
		orderResponse, err = services.SandboxService.PostSandboxOrder(orderRequest)
	} else {
		orderResponse, err = services.OrdersService.PostOrder(orderRequest)
	}
	if err != nil {
		tb.logger.Errorf("can not post sell order: %v", err)
		return // nothing bad happened, let's proceed
	}

	var order Order

	order.SellFlag = true
	order.OrderID = orderResponse.OrderId
	order.OrderPrice = &pb.MoneyValue{
		Units:    fairPrice.Units,
		Nano:     fairPrice.Nano,
		Currency: orderResponse.InitialOrderPrice.Currency,
	}

	t := time.Now().Unix()
	order.OrderPlacedTime = &t

	tb.logger.With("order_id", order.OrderID).
		Infof("sell order created, fair price: %d.%d, initial price: %d.%d %s, current status: %s",
			fairPrice.Units, fairPrice.Nano,
			orderResponse.InitialOrderPrice.Units, orderResponse.InitialOrderPrice.Nano,
			orderResponse.InitialOrderPrice.Currency, orderResponse.ExecutionReportStatus.String())

	metrics.OrdersPlaced.WithLabelValues(loggy.GetBotID(), orderBook.Figi,
		pb.OrderDirection_ORDER_DIRECTION_SELL.String()).Inc()

	tb.orders[orderBook.Figi] = order
}

// checkPortfolio calls sdk.OperationsService.GetPortfolio and updates portfolio metrics.
func (tb *TradeBot) checkPortfolio() {
	var portfolio *pb.PortfolioResponse
	var err error

	if config.TradeBotConfig().IsSandbox {
		portfolio, err = services.SandboxService.GetSandboxPortfolio(tb.accountID)
	} else {
		portfolio, err = services.OperationsService.GetPortfolio(tb.accountID)
	}

	if err != nil {
		tb.logger.Errorf("error getting order book: %v", err)
		return // just ignoring it
	}

	tb.logger.Info("positions: ", tradeutil.GetFormattedPositions(portfolio.Positions))
	common.SetPortfolioMetrics(*portfolio, tb.accountID)

	if portfolio.ExpectedYield != nil {
		tb.logger.Infof("expected yield: %d.%d", portfolio.ExpectedYield.Units, portfolio.ExpectedYield.Nano)
		metrics.PortfolioExpectedYieldOverall.WithLabelValues(tb.accountID).Set(tradeutil.QuotationToFloat(*portfolio.ExpectedYield))
	}
}
