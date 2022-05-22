package main

import (
	"flag"
	"fmt"
	"github.com/elkopass/BITA/internal/config"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/elkopass/BITA/internal/sdk"
	tradeutil "github.com/elkopass/BITA/internal/trade/util"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"time"
)

var services *sdk.ServicePool

func main() {
	services = sdk.NewServicePool()

	var mode string
	flag.StringVar(&mode, "mode", "", "running module")
	flag.Parse()

	if len(mode) == 0 {
		fmt.Println("Usage: trade-utils -mode [accounts|figi|operations]")
		flag.PrintDefaults()
		os.Exit(1)
	}
	switch mode {
	case "accounts":
		printAvailableAccounts()
	case "figi":
		printAvailableFigiList()
	case "operations":
		printLastOperations()
	default:
		fmt.Printf("unknown mode '%s'; possible values: accounts, figi", mode)
		os.Exit(1)
	}
}

// printAvailableAccounts gets accounts from sdk.UsersService.GetAccounts.
func printAvailableAccounts() {
	accounts, err := services.UsersService.GetAccounts()
	if err != nil {
		fmt.Printf("error getting accounts: %v", err)
		os.Exit(1)
	}

	fmt.Println("Available accounts:")
	for _, acc := range accounts {
		fmt.Printf("[%s] %s (%s, %s)\n", acc.Id, acc.Name, acc.Status, acc.AccessLevel)

		positions, err := services.OperationsService.GetPositions(acc.Id)
		if err != nil {
			fmt.Printf("not enough rights to get portfolio: %v \n", err)
		} else {
			fmt.Printf("Funds in the account:\n")
			for _, mon := range positions.Money {
				fmt.Printf("%s: %d.%d\n", mon.Currency, mon.Units, mon.Nano)
			}
			portfolio, err := services.OperationsService.GetPortfolio(acc.Id)
			if err != nil {
				fmt.Printf("not enough rights to get portfolio: %v \n", err)
			} else {
				printPortfolio(*portfolio)
			}
		}
	}
}

// printPortfolio gets instruments information from portfolio positions.
func printPortfolio(portfolio pb.PortfolioResponse) {
	fmt.Printf("Available tools:\n")
	for _, pos := range portfolio.Positions {
		instrument, err := services.InstrumentsService.GetInstrumentBy(
			pb.InstrumentRequest{Id: pos.Figi, IdType: pb.InstrumentIdType_INSTRUMENT_ID_TYPE_FIGI})
		if err != nil {
			fmt.Printf("error getting asset %s: %v\n", pos.Figi, err)
			continue
		}

		candles, err := services.MarketDataService.GetCandles(
			pos.Figi,
			timestamppb.New(time.Now().Add(-24*7*time.Hour)),
			timestamppb.Now(),
			pb.CandleInterval_CANDLE_INTERVAL_HOUR,
		)
		if err != nil {
			fmt.Printf("error getting candles for %s: %v\n", pos.Figi, err)
			continue
		}

		volume, liquidity := tradeutil.CalculateVolumeAndLiquidity(candles)
		averagePrice := tradeutil.MoneyValueToFloat(*pos.AveragePositionPrice)
		currentPrice := tradeutil.MoneyValueToFloat(*pos.CurrentPrice)
		currency := pos.AveragePositionPrice.Currency
		yield := (currentPrice / averagePrice - 1) * 100

		fmt.Printf("[%s] %s (%s, %s): %d quantity, %.2f %s current price (average is %.2f %s), %.2f%% yield, %d volume, %d liquidity\n",
			pos.Figi, instrument.Name, instrument.Ticker, pos.InstrumentType, pos.Quantity.Units,
			currentPrice, currency, averagePrice, currency,
			yield, volume, liquidity)
	}
}

// printAvailableFigiList gets shares and etfs with normal trading status.
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

// printLastOperations gets operations from sdk.OperationsService.GetOperations.
func printLastOperations() {
	services := sdk.NewServicePool()
	operations, err := services.OperationsService.GetOperations(
		config.TradeBotConfig().AccountID,
		timestamppb.New(time.Now().Add(-24*time.Hour)),
		timestamppb.Now(),
		pb.OperationState_OPERATION_STATE_EXECUTED,
		"",
	)
	if err != nil {
		fmt.Printf("error getting operations: %v", err)
		os.Exit(1)
	}

	totalIncome := make(map[string]pb.MoneyValue)

	fmt.Println("Executed orders (last 24 hours):")
	for _, o := range operations {
		mt := totalIncome[o.Currency]

		switch o.OperationType {
		case pb.OperationType_OPERATION_TYPE_SELL:
			mt.Units += o.Price.Units
			mt.Nano += o.Price.Nano
		case pb.OperationType_OPERATION_TYPE_BUY:
			mt.Units -= o.Price.Units
			mt.Nano -= o.Price.Nano
		case pb.OperationType_OPERATION_TYPE_BROKER_FEE:
			mt.Units -= o.Price.Units
			mt.Nano -= o.Price.Nano
		default:
			fmt.Printf("%s is not supported!\n", o.OperationType.String())
		}
		totalIncome[o.Currency] = mt

		fmt.Printf("[%s] %s: %d.%d, %s (%s)\n",
			o.Figi, o.OperationType.String(), o.Price.Units, o.Price.Nano, o.Currency, o.Date.AsTime())
	}

	fmt.Println()
	fmt.Printf("total income:\n")
	for currency, mt := range totalIncome {
		fmt.Printf("%d.%d %s\n", mt.Units, mt.Nano, currency)
	}
}
