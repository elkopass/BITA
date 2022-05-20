package main

import (
	"flag"
	"fmt"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/elkopass/BITA/internal/sdk"
	"os"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "", "running module")
	flag.Parse()

	if len(mode) == 0 {
		fmt.Println("Usage: trade-utils -mode [accounts|figi]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch mode {
	case "accounts":
		printAvailableAccounts()
	case "figi":
		printAvailableFigiList()
	default:
		fmt.Printf("unknown mode '%s'; possible values: accounts, figi", mode)
		os.Exit(1)
	}
}

func printAvailableFigiList() {
	services := sdk.NewServicePool()
	shares, err := services.InstrumentsService.Shares(pb.InstrumentStatus_INSTRUMENT_STATUS_ALL)
	if err != nil {
		fmt.Printf("error getting shares: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Available shares:")
	for _, share := range shares {
		if share.TradingStatus == pb.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING {
			message := fmt.Sprintf("[%s] %s: %s", share.Ticker, share.Name, share.Figi)
			message += fmt.Sprintf(" (currency: %s, lot: %d)", share.Currency, share.Lot)

			fmt.Println(message)
		}
	}

	etfs, err := services.InstrumentsService.Etfs(pb.InstrumentStatus_INSTRUMENT_STATUS_ALL)
	if err != nil {
		fmt.Printf("error getting etfs: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Available efts:")
	for _, etf := range etfs {
		if etf.TradingStatus == pb.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING {
			message := fmt.Sprintf("[%s] %s: %s", etf.Ticker, etf.Name, etf.Figi)
			message += fmt.Sprintf(" (currency: %s, lot: %d)", etf.Currency, etf.Lot)

			fmt.Println(message)
		}
	}
}

func printAvailableAccounts() {
	services := sdk.NewServicePool()
	accounts, err := services.UsersService.GetAccounts()
	if err != nil {
		fmt.Printf("error getting accounts: %v", err)
		os.Exit(1)
	}

	fmt.Println("Available accounts:")
	for _, acc := range accounts {
		fmt.Printf("[%s] %s (%s, %s)\n", acc.Id, acc.Name, acc.Status, acc.AccessLevel)
	}
}
