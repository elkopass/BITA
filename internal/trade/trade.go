package trade

import "context"

type Trader interface {
	Run(ctx context.Context) (err error)
}
