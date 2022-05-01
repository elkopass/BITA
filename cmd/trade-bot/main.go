package main

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/config"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
)


func main() {
	log := loggy.GetLogger().Sugar()

	tradeBotConfig := config.TradeBotConfig()
	if tradeBotConfig.Token == "" {
		log.Fatal("TRADEBOT_TOKEN environment variable is required to run this program")
	}
}
