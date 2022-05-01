package sdk

import (
	"github.com/elkopass/TinkoffInvestRobotContest/internal/loggy"
	pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// TODO: implementation
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
	GetAssetBy(filters InstrumentSearchFilters) (*pb.AssetFull, error)
	// Метод получения списка активов.
	GetAssets() ([]*pb.Asset, error)
	// Метод получения избранных инструментов.
	GetFavorites() ([]*pb.FavoriteInstrument, error)
	// Метод редактирования избранных инструментов.
	EditFavorites(newFavourites *pb.EditFavoritesRequest) ([]*pb.FavoriteInstrument, error)
}

type InstrumentsService struct {
	client *pb.InstrumentsServiceClient
}

func NewInstrumentsService() *InstrumentsService {
	conn, err := createClientConn()
	if err != nil {
		loggy.GetLogger().Sugar().Fatal(err.Error())
	}

	client := pb.NewInstrumentsServiceClient(conn)
	return &InstrumentsService{client: &client}
}
