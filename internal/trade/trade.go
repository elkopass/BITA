// Package trade is responsible for all trading logic.
package trade

import (
	"context"
	"sync"
)

// Trader is a common interface for all trading bots.
type Trader interface {
	// Run starts a trading bot with multiple workers;
	// context.Context is required for cancellation.
	Run(ctx context.Context) (err error)
}

// TradeWorker is a common interface for all trade workers.
type TradeWorker interface {
	// Run starts a worker; context.Context is required for cancellation
	// and a *sync.WaitGroup is required for graceful shutdown.
	Run(ctx context.Context, wg *sync.WaitGroup) (err error)
}
