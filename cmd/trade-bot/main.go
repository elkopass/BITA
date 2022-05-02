package main

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/config"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/trade"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/trade/strategy"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/trade/strategy/gamble"
)

func main() {
	log := loggy.GetLogger().Sugar()

	cnf := config.TradeBotConfig()
	if cnf.IsSandbox {
		log.Infof("running in sandbox mode with %s strategy", cnf.Strategy)
	} else {
		log.Warnf("[DANGER] running without sandbox with %s strategy, I hope you know what you doing", cnf.Strategy)
	}

	var bot trade.Trader
	switch cnf.Strategy {
	case strategy.GAMBLE:
		bot = gamble.NewTraderBot()
	default:
		log.Fatalf("unknown strategy '%s'", cnf.Strategy)
		return
	}

	if cnf.IsSandbox {
		bot.RunInSandbox()
	} else {
		bot.Run()
	}
}
