// Package trade is responsible for all trading logic.
package trade

import "context"

type Trader interface {
	Run(ctx context.Context) (err error)
}
