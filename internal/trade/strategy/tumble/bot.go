package tumble

import (
	"context"
	"fmt"
	"github.com/elkopass/BITA/internal/config"
	"github.com/elkopass/BITA/internal/loggy"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/elkopass/BITA/internal/sdk"
	tradeutil "github.com/elkopass/BITA/internal/trade/util"
	"go.uber.org/zap"
	"strings"
	"time"
)

type TradeBot struct {
	config      TradeConfig
	logger      *zap.SugaredLogger
}

func NewTradeBot() *TradeBot {
	return &TradeBot{
		config: *NewTradeConfig(),
		logger: loggy.GetLogger().Sugar().With("bot_id", loggy.GetBotID()),
	}
}

func (tb TradeBot) Run(ctx context.Context) (err error) {
	tb.logger.Infof("starting with %s strategy and sdk v%s", config.TradeBotConfig().Strategy, sdk.Version)

	accountID := config.TradeBotConfig().AccountID
	if config.TradeBotConfig().IsSandbox {
		accountID, err = services.SandboxService.OpenSandboxAccount()
		if err != nil {
			return fmt.Errorf("can not create account: %v", err)
		}
		tb.logger.Infof("created new account with ID %s", accountID)
	} else {
		info, err := services.UsersService.GetInfo()
		if err != nil {
			return fmt.Errorf("can not get user info: %v", err)
		}
		tb.logger.Infof("user tariff: %s, qualified for work with %s",
			info.Tariff, strings.Join(info.QualifiedForWorkWith, ","))
	}

	// replace logger
	tb.logger = tb.logger.With("account_id", accountID)

	figi := config.TradeBotConfig().Figi

	mdss := sdk.NewMarketDataStreamService()
	stream, err := mdss.MarketDataStream()
	if err != nil {
		panic(err)
	}

	var instruments []*pb.OrderBookInstrument
	for _, f := range figi {
		instruments = append(instruments, &pb.OrderBookInstrument{Figi: f, Depth: 10})
	}

	request := pb.SubscribeOrderBookRequest{
		Instruments: instruments,
		SubscriptionAction: pb.SubscriptionAction_SUBSCRIPTION_ACTION_SUBSCRIBE,
	}
	payload := &pb.MarketDataRequest_SubscribeOrderBookRequest{SubscribeOrderBookRequest: &request}

	err = stream.Send(&pb.MarketDataRequest{Payload: payload})
	if err != nil {
		return err
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			tb.logger.Error(err)
		}

		orderBook := res.GetOrderbook()
		tb.logger.Info(tradeutil.GetFormattedOrderBook(orderBook))

		select {
		case <-time.After(1 * time.Millisecond):
			// pass
		case <-ctx.Done():
			// TODO: implement sell logic on interrupt
			if config.TradeBotConfig().IsSandbox {
				err = services.SandboxService.CloseSandboxAccount(accountID)
				if err != nil {
					tb.logger.Errorf("can't create account: %v", err)
				}
				tb.logger.Infof("account with ID %s closed successfully", accountID)
			}

			tb.logger.Info("bot stopped!")
			return nil
		}
	}
}
