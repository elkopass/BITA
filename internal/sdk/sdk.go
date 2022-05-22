// Package sdk represents internal proto-wrapper for Tinkoff Invest API.
package sdk

const Version = "0.2.1"

// ServicePool is a ready-to-use scope for all available non-stream services.
type ServicePool struct {
	InstrumentsService InstrumentsService
	MarketDataService  MarketDataService
	OperationsService  OperationsService
	OrdersService      OrdersService
	SandboxService     SandboxService
	StopOrdersService  StopOrdersService
	UsersService       UsersService
}

func NewServicePool() *ServicePool {
	return &ServicePool{
		InstrumentsService: *NewInstrumentsService(),
		MarketDataService:  *NewMarketDataService(),
		OperationsService:  *NewOperationsService(),
		OrdersService:      *NewOrdersService(),
		StopOrdersService:  *NewStopOrdersService(),
		SandboxService:     *NewSandboxService(),
		UsersService:       *NewUsersService(),
	}
}
