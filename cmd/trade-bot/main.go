package main

import (
	"context"
	"github.com/elkopass/BITA/internal/config"
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/trade"
	"github.com/elkopass/BITA/internal/trade/strategy"
	"github.com/elkopass/BITA/internal/trade/strategy/gamble"
	"github.com/elkopass/BITA/internal/trade/strategy/tumble"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	log := loggy.GetLogger().Sugar()

	// validate configuration
	cnf := config.TradeBotConfig()
	if cnf.Token == "<your_api_token>" {
		log.Fatalf("please set your own API token in TRADEBOT_TOKEN env variable")
	}

	if cnf.IsSandbox {
		log.Infof("running in sandbox mode with %s strategy", cnf.Strategy)
	} else {
		if cnf.AccountID == "<your_api_token>" {
			log.Fatalf("please specify your own account ID in TRADEBOT_ACCOUNT_ID env variable " +
				"(compile and run 'trade-utils -mode accounts' to get it)")
		}
		log.Warnf("[DANGER] running without sandbox with %s strategy and %s account ID, " +
			"I hope you know what you doing", cnf.Strategy, cnf.AccountID)
	}

	if len(cnf.Figi) == 2 && cnf.Figi[0] == "<figi1>" {
		log.Fatalf("please specify some figi's to trade in TRADEBOT_FIGI env variable; " +
			"if you need some, compile and run 'trade-utils -mode figi' to get them")
	}

	// init trade bot
	var bot trade.Trader

	switch cnf.Strategy {
	case strategy.GAMBLE:
		bot = gamble.NewTradeBot()
	case strategy.TUMBLE:
		bot = tumble.NewTradeBot()
	default:
		log.Fatalf("unknown strategy '%s'", cnf.Strategy)
		return
	}

	// setting up server for metrics
	metrics := config.MetricsConfig()
	if metrics.Enabled {
		server := http.NewServeMux()
		server.Handle(metrics.Endpoint, promhttp.Handler())

		srv := &http.Server{
			Addr:    metrics.Addr,
			Handler: server,
		}

		go func() {
			log.Infof("listening on %s", metrics.Addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen failed: %+s\n", err)
			}
		}()
	}

	// preparing for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		s := <-c
		log.Warnf("received system call: %+v", s)
		cancel()
	}()

	// run bot until interrupt signal is received
	if err := bot.Run(ctx); err != nil {
		log.Errorf("failed to shutdown trade bot: +%v\n", err)
	}

	log.Info("trade bot exited properly")
}
