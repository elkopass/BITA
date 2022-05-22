// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

import (
	"github.com/elkopass/BITA/internal/loggy"
	"github.com/elkopass/BITA/internal/metrics"
	pb "github.com/elkopass/BITA/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type InstrumentsInterface interface {
	// The method of obtaining the trading schedule of trading platforms.
	TradingSchedules(exchange string, from, to *timestamp.Timestamp) ([]*pb.TradingSchedule, error)
	// The method of obtaining a bond by its identifier.
	BondBy(filters pb.InstrumentRequest) (*pb.Bond, error)
	// Method of obtaining a list of bonds.
	Bonds(status pb.InstrumentStatus) ([]*pb.Bond, error)
	// Method of obtaining a coupon payment schedule for a bond.
	GetBondCoupons(figi string, from, to *timestamp.Timestamp) ([]*pb.Coupon, error)
	// The method of obtaining a currency by its identifier.
	CurrencyBy(filters pb.InstrumentRequest) (*pb.Currency, error)
	// Method for getting a list of currencies.
	Currencies(status pb.InstrumentStatus) ([]*pb.Currency, error)
	// The method of obtaining an investment fund by its identifier.
	EtfBy(filters pb.InstrumentRequest) (*pb.Etf, error)
	// Method of obtaining a list of investment funds.
	Etfs(status pb.InstrumentStatus) ([]*pb.Etf, error)
	// The method of obtaining futures by its identifier.
	FutureBy(filters pb.InstrumentRequest) (*pb.Future, error)
	// Method for getting a list of futures.
	Futures(status pb.InstrumentStatus) ([]*pb.Future, error)
	// The method of obtaining a stock by its identifier.
	ShareBy(filters pb.InstrumentRequest) (*pb.Share, error)
	// Method of getting a list of shares.
	Shares(status pb.InstrumentStatus) ([]*pb.Share, error)
	// The method of obtaining the accumulated coupon income on the bond.
	GetAccruedInterests(figi string, from, to *timestamp.Timestamp) ([]*pb.AccruedInterest, error)
	// The method of obtaining the amount of the guarantee for futures.
	GetFuturesMargin(figi string) (*pb.GetFuturesMarginResponse, error)
	// The method of obtaining basic information about the tool.
	GetInstrumentBy(filters pb.InstrumentRequest) (*pb.Instrument, error)
	// A method for obtaining dividend payment events for an instrument.
	GetDividends(figi string, from, to *timestamp.Timestamp) ([]*pb.Dividend, error)
	// The method of obtaining an asset by its identifier.
	GetAssetBy(assetID string) (*pb.AssetFull, error)
	// Method for getting a list of assets.
	GetAssets() ([]*pb.Asset, error)
	// The method of getting the favourite instruments.
	GetFavorites() ([]*pb.FavoriteInstrument, error)
	// The method of editing the selected instruments.
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
		is.incrementApiCallErrors("TradingSchedules", err.Error())
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
		is.incrementApiCallErrors("BondBy", err.Error())
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
		is.incrementApiCallErrors("Bonds", err.Error())
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
		is.incrementApiCallErrors("GetBoundCoupons", err.Error())
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
		is.incrementApiCallErrors("CurrencyBy", err.Error())
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
		is.incrementApiCallErrors("Currencies", err.Error())
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
		is.incrementApiCallErrors("EtfBy", err.Error())
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
		is.incrementApiCallErrors("Etfs", err.Error())
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
		is.incrementApiCallErrors("FutureBy", err.Error())
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
		is.incrementApiCallErrors("Futures", err.Error())
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
		is.incrementApiCallErrors("ShareBy", err.Error())
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
		is.incrementApiCallErrors("Shares", err.Error())
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
		is.incrementApiCallErrors("GetAccruedInterests", err.Error())
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
		is.incrementApiCallErrors("GetFuturesMargin", err.Error())
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
		is.incrementApiCallErrors("GetInstrumentBy", err.Error())
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
		is.incrementApiCallErrors("GetDividends", err.Error())
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
		is.incrementApiCallErrors("GetAssetBy", err.Error())
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
		is.incrementApiCallErrors("GetAssets", err.Error())
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
		is.incrementApiCallErrors("GetFavourites", err.Error())
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
		is.incrementApiCallErrors("EditFavorites", err.Error())
		return nil, err
	}

	return res.FavoriteInstruments, nil
}

func (is InstrumentsService) incrementRequestsCounter(method string) {
	metrics.ApiRequests.WithLabelValues(loggy.GetBotID(), "InstrumentsService", method).Inc()
}

func (is InstrumentsService) incrementApiCallErrors(method string, error string) {
	metrics.ApiCallErrors.WithLabelValues(loggy.GetBotID(), "InstrumentsService", method, error).Inc()
}
