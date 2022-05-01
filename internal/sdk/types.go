package sdk

import pb "github.com/elkopass/TinkoffInvestRobotContest/internal/proto"

type Figi string
type AccountID string
type OrderID string
type StopOrderID string

type InstrumentSearchFilters pb.InstrumentRequest
type OperationsSearchFilters pb.OperationsRequest

type Order pb.PostOrderRequest
type OrderState pb.OrderState
type PostOrder pb.PostOrderResponse

type Operation pb.Operation
type Portfolio pb.PortfolioResponse
type Positions pb.PositionsResponse
type WithdrawLimits pb.WithdrawLimitsResponse

type TradingStatus pb.GetTradingStatusResponse
type OrderBook pb.GetOrderBookResponse
