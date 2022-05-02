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

	tradeBotConfig := config.TradeBotConfig()

	var bot trade.Trader
	switch tradeBotConfig.Strategy {
		case strategy.GAMBLE:
			bot = gamble.NewTraderBot()
		default:
			log.Fatalf("unknown strategy '%s'", tradeBotConfig.Strategy)
			return
	}

	bot.Run()
}
