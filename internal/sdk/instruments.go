package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	"github.com/elkopass/TinkoffInvestRobotContest/internal/metrics"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type InstrumentsInterface interface {
	// Метод получения расписания торгов торговых площадок.
	TradingSchedules(exchange string, from, to *timestamp.Timestamp) ([]*pb.TradingSchedule, error)
	// Метод получения облигации по её идентификатору.
	BondBy(filters pb.InstrumentRequest) (*pb.Bond, error)
	// Метод получения списка облигаций.
	Bonds(status pb.InstrumentStatus) ([]*pb.Bond, error)
	// Метод получения графика выплат купонов по облигации
	GetBondCoupons(figi string, from, to *timestamp.Timestamp) ([]*pb.Coupon, error)
	// Метод получения валюты по её идентификатору.
	CurrencyBy(filters pb.InstrumentRequest) (*pb.Currency, error)
	// Метод получения списка валют.
	Currencies(status pb.InstrumentStatus) ([]*pb.Currency, error)
	// Метод получения инвестиционного фонда по его идентификатору.
	EtfBy(filters pb.InstrumentRequest) (*pb.Etf, error)
	// Метод получения списка инвестиционных фондов.
	Etfs(status pb.InstrumentStatus) ([]*pb.Etf, error)
	// Метод получения фьючерса по его идентификатору.
	FutureBy(filters pb.InstrumentRequest) (*pb.Future, error)
	// Метод получения списка фьючерсов.
	Futures(status pb.InstrumentStatus) ([]*pb.Future, error)
	// Метод получения акции по её идентификатору.
	ShareBy(filters pb.InstrumentRequest) (*pb.Share, error)
	// Метод получения списка акций.
	Shares(status pb.InstrumentStatus) ([]*pb.Share, error)
	// Метод получения накопленного купонного дохода по облигации.
	GetAccruedInterests(figi string, from, to *timestamp.Timestamp) ([]*pb.AccruedInterest, error)
	// Метод получения размера гарантийного обеспечения по фьючерсам.
	GetFuturesMargin(figi string) (*pb.GetFuturesMarginResponse, error)
	// Метод получения основной информации об инструменте.
	GetInstrumentBy(filters pb.InstrumentRequest) (*pb.Instrument, error)
	// Метод для получения событий выплаты дивидендов по инструменту.
	GetDividends(figi string, from, to *timestamp.Timestamp) ([]*pb.Dividend, error)
	// Метод получения актива по его идентификатору.
	GetAssetBy(assetID string) (*pb.AssetFull, error)
	// Метод получения списка активов.
	GetAssets() ([]*pb.Asset, error)
	// Метод получения избранных инструментов.
	GetFavorites() ([]*pb.FavoriteInstrument, error)
	// Метод редактирования избранных инструментов.
	EditFavorites(newFavourites *pb.EditFavoritesRequest) ([]*pb.FavoriteInstrument, error)
}

type InstrumentsService struct {
	client pb.InstrumentsServiceClient
}

func NewInstrumentsService() *InstrumentsService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewInstrumentsServiceClient(conn)
	return &InstrumentsService{client: client}
}

func (is InstrumentsService) TradingSchedules(exchange string, from, to *timestamp.Timestamp) ([]*pb.TradingSchedule, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("TradingSchedules")
	res, err := is.client.TradingSchedules(ctx, &pb.TradingSchedulesRequest{
		Exchange: exchange,
		From:     from,
		To:       to,
	})
	if err != nil {
		return nil, err
	}

	return res.Exchanges, nil
}

func (is InstrumentsService) BondBy(filters pb.InstrumentRequest) (*pb.Bond, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("BondBy")
	res, err := is.client.BondBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Bonds(status pb.InstrumentStatus) ([]*pb.Bond, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("Bonds")
	res, err := is.client.Bonds(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) GetBondCoupons(figi string, from, to *timestamp.Timestamp) ([]*pb.Coupon, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("GetBoundCoupons")
	res, err := is.client.GetBondCoupons(ctx, &pb.GetBondCouponsRequest{
		Figi: figi,
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, err
	}

	return res.Events, nil
}

func (is InstrumentsService) CurrencyBy(filters pb.InstrumentRequest) (*pb.Currency, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("CurrencyBy")
	res, err := is.client.CurrencyBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Currencies(status pb.InstrumentStatus) ([]*pb.Currency, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("Currencies")
	res, err := is.client.Currencies(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) EtfBy(filters pb.InstrumentRequest) (*pb.Etf, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("EtfBy")
	res, err := is.client.EtfBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Etfs(status pb.InstrumentStatus) ([]*pb.Etf, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("Etfs")
	res, err := is.client.Etfs(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) FutureBy(filters pb.InstrumentRequest) (*pb.Future, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("FutureBy")
	res, err := is.client.FutureBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Futures(status pb.InstrumentStatus) ([]*pb.Future, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("Futures")
	res, err := is.client.Futures(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) ShareBy(filters pb.InstrumentRequest) (*pb.Share, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("ShareBy")
	res, err := is.client.ShareBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Shares(status pb.InstrumentStatus) ([]*pb.Share, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("Shares")
	res, err := is.client.Shares(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) GetAccruedInterests(figi string, from, to *timestamp.Timestamp) ([]*pb.AccruedInterest, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("GetAccruedInterests")
	res, err := is.client.GetAccruedInterests(ctx, &pb.GetAccruedInterestsRequest{
		Figi: figi,
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, err
	}

	return res.AccruedInterests, nil
}

func (is InstrumentsService) GetFuturesMargin(figi string) (*pb.GetFuturesMarginResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("GetFuturesMargin")
	res, err := is.client.GetFuturesMargin(ctx, &pb.GetFuturesMarginRequest{
		Figi: figi,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (is InstrumentsService) GetInstrumentBy(filters pb.InstrumentRequest) (*pb.Instrument, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("GetInstrumentBy")
	res, err := is.client.GetInstrumentBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) GetDividends(figi string, from, to *timestamp.Timestamp) ([]*pb.Dividend, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("GetDividends")
	res, err := is.client.GetDividends(ctx, &pb.GetDividendsRequest{
		Figi: figi,
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, err
	}

	return res.Dividends, nil
}

func (is InstrumentsService) GetAssetBy(assetID string) (*pb.AssetFull, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("GetAssetBy")
	res, err := is.client.GetAssetBy(ctx, &pb.AssetRequest{
		Id: assetID,
	})
	if err != nil {
		return nil, err
	}

	return res.Asset, nil
}

func (is InstrumentsService) GetAssets() ([]*pb.Asset, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("GetAssets")
	res, err := is.client.GetAssets(ctx, &pb.AssetsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Assets, nil
}

func (is InstrumentsService) GetFavorites() ([]*pb.FavoriteInstrument, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("GetFavourites")
	res, err := is.client.GetFavorites(ctx, &pb.GetFavoritesRequest{})
	if err != nil {
		return nil, err
	}

	return res.FavoriteInstruments, nil
}

func (is InstrumentsService) EditFavorites(newFavourites *pb.EditFavoritesRequest) ([]*pb.FavoriteInstrument, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	is.incrementRequestsCounter("EditFavorites")
	res, err := is.client.EditFavorites(ctx, newFavourites)
	if err != nil {
		return nil, err
	}

	return res.FavoriteInstruments, nil
}

func (is InstrumentsService) incrementRequestsCounter(method string) {
	metrics.ApiRequests.WithLabelValues(loggy.GetBotID(), "InstrumentsService", method).Inc()
}
