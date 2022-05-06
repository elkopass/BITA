package main

import (
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
			log.Infof("%s: %s (%s)", share.Name, share.Figi, share.TradingStatus)
		}
	}
}
