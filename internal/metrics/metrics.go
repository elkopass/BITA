package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	InstrumentsPurchased = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_instruments_purchased",
		Help: "Purchased instruments gauge",
	}, []string{"figi", "bot_id"})
	OrdersPlaced = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tradebot_orders_placed",
		Help: "Placed orders gauge",
	}, []string{"figi", "bot_id"})
	OrdersFulfilled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tradebot_orders_fulfilled",
		Help: "Fulfilled orders counter",
	}, []string{"figi", "bot_id"})
)

func init()  {
	prometheus.MustRegister(InstrumentsPurchased)
	prometheus.MustRegister(OrdersPlaced)
	prometheus.MustRegister(OrdersFulfilled)
}
