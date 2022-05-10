package sdk

const Version = "0.1.4"

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
