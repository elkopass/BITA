package main

import (
	"fmt"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/sdk"
)

func main() {
	log := loggy.GetLogger().Sugar()

	services := sdk.NewServicePool()
	shares, err := services.InstrumentsService.Shares(pb.InstrumentStatus_INSTRUMENT_STATUS_ALL)
	if err != nil {
		log.Fatalf("error getting shares: %v", err)
		return
	}

	for _, share := range shares {
		if share.TradingStatus == pb.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING {
			message := fmt.Sprintf("[%s] %s: %s", share.Ticker, share.Name, share.Figi)
			message += fmt.Sprintf(" (currency: %s, lot: %d)", share.Currency, share.Lot)

			log.Info(message)
		}
	}
}
