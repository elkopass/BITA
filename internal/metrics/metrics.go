// Package metrics stores all Prometheus metrics used by trade-bot.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// BotInfo stores strategy, sdk version and amount of workers for each running bot.
	BotInfo = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_info",
		Help: "Brief information about running bot",
	}, []string{"bot_id", "sdk_version", "selected_strategy", "figi_amount"})

	// ApiRequests counts total requests to Tinkoff Invest API.
	ApiRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_api_requests",
		Help: "Total requests to Tinkoff Invest API counter",
	}, []string{"bot_id", "service", "method"})
	// ApiCallErrors counts total number of failed API requests.
	ApiCallErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_api_call_errors",
		Help: "Total failed requests to Tinkoff Invest API counter",
	}, []string{"bot_id", "service", "method", "error"})

	// InstrumentsPurchased stores number of currently purchased instruments.
	InstrumentsPurchased = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_instruments_purchased",
		Help: "Purchased instruments by bot gauge",
	}, []string{"bot_id", "figi"})
	// OrdersPlaced stores number of currently places buy/sell orders.
	OrdersPlaced = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_orders_placed",
		Help: "Placed orders by bot gauge",
	}, []string{"bot_id", "figi", "direction"})
	// OrdersFulfilled counts number of total sold instruments.
	OrdersFulfilled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_orders_fulfilled",
		Help: "Fulfilled orders total counter",
	}, []string{"bot_id", "figi", "direction"})
	// OrdersCancelled counts number or cancelled orders.
	OrdersCancelled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_orders_cancelled",
		Help: "Cancelled orders total counter",
	}, []string{"bot_id", "figi"})
	// StopLossDecisions counts number of stop loss decisions by trade bot.
	StopLossDecisions = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_stop_loss_decisions",
		Help: "Stop loss decisions counter",
	}, []string{"bot_id", "figi"})
	// TakeProfitDecisions counts number of stop loss decisions by trade bot.
	TakeProfitDecisions = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_take_profit_decisions",
		Help: "Take profit decisions counter",
	}, []string{"bot_id", "figi"})
	// StoppedByCircuitBreaker counts unhealthy workers removed by circuit breaker.
	StoppedByCircuitBreaker = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_stopped_by_circuit_breaker",
		Help: "Trade workers stopped by circuit breaker counter",
	}, []string{"bot_id", "figi"})

	// InstrumentLastPrice stores last price for existing instrument.
	InstrumentLastPrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_instrument_last_price",
		Help: "Instrument last price gauge",
	}, []string{"figi"})
	// InstrumentFairPrice stores calculated fair price for an instrument.
	InstrumentFairPrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_instrument_fair_price",
		Help: "Instrument fair price gauge",
	}, []string{"figi"})
	// InstrumentTradingStatus stores last instrument trading status.
	InstrumentTradingStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_instrument_trading_status",
		Help: "Instrument trading status gauge",
	}, []string{"figi", "status"})
	// PortfolioInstrumentsAmount stores portfolio instruments amount by each type.
	PortfolioInstrumentsAmount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_portfolio_instruments_amount",
		Help: "Portfolio instruments total money amount gauge",
	}, []string{"account_id", "instrument", "currency"})
	// PortfolioExpectedYieldOverall stores expected portfolio income.
	PortfolioExpectedYieldOverall = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_portfolio_expected_yield_overall",
		Help: "Portfolio expected yield for all positions gauge",
	}, []string{"account_id"})
	// PortfolioPositionCurrentPrice stores current price for each open position.
	PortfolioPositionCurrentPrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_portfolio_position_current_price",
		Help: "Portfolio position current price (for one lot) gauge",
	}, []string{"account_id", "figi"})
	// PortfolioPositionCurrentPrice stores expected yield for each open position.
	PortfolioPositionExpectedYield = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_portfolio_position_expected_yield",
		Help: "Portfolio position expected yield gauge",
	}, []string{"account_id", "figi"})
)

func init() {
	prometheus.MustRegister(BotInfo)

	/* sdk related metrics */
	prometheus.MustRegister(ApiRequests)
	prometheus.MustRegister(ApiCallErrors)

	/* bot key actions */
	prometheus.MustRegister(InstrumentsPurchased)
	prometheus.MustRegister(OrdersPlaced)
	prometheus.MustRegister(OrdersFulfilled)
	prometheus.MustRegister(OrdersCancelled)
	prometheus.MustRegister(StopLossDecisions)
	prometheus.MustRegister(TakeProfitDecisions)
	prometheus.MustRegister(StoppedByCircuitBreaker)

	/* additional trade statistics */
	prometheus.MustRegister(InstrumentLastPrice)
	prometheus.MustRegister(InstrumentFairPrice)
	prometheus.MustRegister(InstrumentTradingStatus)
	prometheus.MustRegister(PortfolioInstrumentsAmount)
	prometheus.MustRegister(PortfolioExpectedYieldOverall)
	prometheus.MustRegister(PortfolioPositionCurrentPrice)
	prometheus.MustRegister(PortfolioPositionExpectedYield)
}
