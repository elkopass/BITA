package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type InstrumentsInterface interface {
	// Метод получения расписания торгов торговых площадок.
	TradingSchedules(exchange string, from, to *timestamp.Timestamp) ([]*pb.TradingSchedule, error)
	// Метод получения облигации по её идентификатору.
	BondBy(filters InstrumentSearchFilters) (*pb.Bond, error)
	// Метод получения списка облигаций.
	Bonds(status pb.InstrumentStatus) ([]*pb.Bond, error)
	// Метод получения графика выплат купонов по облигации
	GetBondCoupons(figi Figi, from, to *timestamp.Timestamp) ([]*pb.Coupon, error)
	// Метод получения валюты по её идентификатору.
	CurrencyBy(filters InstrumentSearchFilters) (*pb.Currency, error)
	// Метод получения списка валют.
	Currencies(status pb.InstrumentStatus) ([]*pb.Currency, error)
	// Метод получения инвестиционного фонда по его идентификатору.
	EtfBy(filters InstrumentSearchFilters) (*pb.Etf, error)
	// Метод получения списка инвестиционных фондов.
	Etfs(status pb.InstrumentStatus) ([]*pb.Etf, error)
	// Метод получения фьючерса по его идентификатору.
	FutureBy(filters InstrumentSearchFilters) (*pb.Future, error)
	// Метод получения списка фьючерсов.
	Futures(status pb.InstrumentStatus) ([]*pb.Future, error)
	// Метод получения акции по её идентификатору.
	ShareBy(filters InstrumentSearchFilters) (*pb.Share, error)
	// Метод получения списка акций.
	Shares(status pb.InstrumentStatus) ([]*pb.Share, error)
	// Метод получения накопленного купонного дохода по облигации.
	GetAccruedInterests(figi Figi, from, to *timestamp.Timestamp) ([]*pb.AccruedInterest, error)
	// Метод получения размера гарантийного обеспечения по фьючерсам.
	GetFuturesMargin(figi Figi) (*pb.GetFuturesMarginResponse, error)
	// Метод получения основной информации об инструменте.
	GetInstrumentBy(filters InstrumentSearchFilters) (*pb.Instrument, error)
	// Метод для получения событий выплаты дивидендов по инструменту.
	GetDividends(figi Figi, from, to *timestamp.Timestamp) ([]*pb.Dividend, error)
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

func (is InstrumentsService) BondBy(filters InstrumentSearchFilters) (*pb.Bond, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.BondBy(ctx, &pb.InstrumentRequest{
		IdType:    filters.IdType,
		ClassCode: filters.ClassCode,
		Id:        filters.Id,
	})
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Bonds(status pb.InstrumentStatus) ([]*pb.Bond, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.Bonds(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) GetBondCoupons(figi Figi, from, to *timestamp.Timestamp) ([]*pb.Coupon, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.GetBondCoupons(ctx, &pb.GetBondCouponsRequest{
		Figi: string(figi),
		From: from,
		To: to,
	})
	if err != nil {
		return nil, err
	}

	return res.Events, nil
}

func (is InstrumentsService) CurrencyBy(filters InstrumentSearchFilters) (*pb.Currency, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.CurrencyBy(ctx, &pb.InstrumentRequest{
		IdType:    filters.IdType,
		ClassCode: filters.ClassCode,
		Id:        filters.Id,
	})
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Currencies(status pb.InstrumentStatus) ([]*pb.Currency, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.Currencies(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) EtfBy(filters InstrumentSearchFilters) (*pb.Etf, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.EtfBy(ctx, &pb.InstrumentRequest{
		IdType:    filters.IdType,
		ClassCode: filters.ClassCode,
		Id:        filters.Id,
	})
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Etfs(status pb.InstrumentStatus) ([]*pb.Etf, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.Etfs(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) FutureBy(filters InstrumentSearchFilters) (*pb.Future, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.FutureBy(ctx, &pb.InstrumentRequest{
		IdType:    filters.IdType,
		ClassCode: filters.ClassCode,
		Id:        filters.Id,
	})
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Futures(status pb.InstrumentStatus) ([]*pb.Future, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.Futures(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) ShareBy(filters InstrumentSearchFilters) (*pb.Share, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.ShareBy(ctx, &pb.InstrumentRequest{
		IdType:    filters.IdType,
		ClassCode: filters.ClassCode,
		Id:        filters.Id,
	})
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Shares(status pb.InstrumentStatus) ([]*pb.Share, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.Shares(ctx, &pb.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) GetAccruedInterests(figi Figi, from, to *timestamp.Timestamp) ([]*pb.AccruedInterest, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.GetAccruedInterests(ctx, &pb.GetAccruedInterestsRequest{
		Figi: string(figi),
		From: from,
		To: to,
	})
	if err != nil {
		return nil, err
	}

	return res.AccruedInterests, nil
}

func (is InstrumentsService) GetFuturesMargin(figi Figi) (*pb.GetFuturesMarginResponse, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.GetFuturesMargin(ctx, &pb.GetFuturesMarginRequest{
		Figi: string(figi),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (is InstrumentsService) GetInstrumentBy(filters InstrumentSearchFilters) (*pb.Instrument, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.GetInstrumentBy(ctx, &pb.InstrumentRequest{
		IdType:    filters.IdType,
		ClassCode: filters.ClassCode,
		Id:        filters.Id,
	})
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) GetDividends(figi Figi, from, to *timestamp.Timestamp) ([]*pb.Dividend, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.GetDividends(ctx, &pb.GetDividendsRequest{
		Figi: string(figi),
		From: from,
		To: to,
	})
	if err != nil {
		return nil, err
	}

	return res.Dividends, nil
}

func (is InstrumentsService) GetAssetBy(assetID string) (*pb.AssetFull, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

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

	res, err := is.client.GetAssets(ctx, &pb.AssetsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Assets, nil
}

func (is InstrumentsService) GetFavorites() ([]*pb.FavoriteInstrument, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.GetFavorites(ctx, &pb.GetFavoritesRequest{})
	if err != nil {
		return nil, err
	}

	return res.FavoriteInstruments, nil
}

func (is InstrumentsService) EditFavorites(newFavourites *pb.EditFavoritesRequest) ([]*pb.FavoriteInstrument, error) {
	ctx, cancel := createRequestContext()
	defer cancel()

	res, err := is.client.EditFavorites(ctx, newFavourites)
	if err != nil {
		return nil, err
	}

	return res.FavoriteInstruments, nil
}
