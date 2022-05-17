// Package breaker stops trade worker if it violates its failure threshold.
package breaker // import cb "github.com/elkopass/BITA/internal/trade/breaker"

import (
	"github.com/elkopass/BITA/internal/config"
	"time"
)

type CircuitBreaker struct {
	failuresTotal   int
	lastFailureTime *time.Time
}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{}
}

func (cb *CircuitBreaker) AddFailure() {
	cb.failuresTotal++
	cb.updateState()
}

// WorkerMustExit returns true if trade worker is unhealthy and must be killed.
func (cb CircuitBreaker) WorkerMustExit() bool {
	return cb.failuresTotal > config.CircuitBreakerMaxFailures
}

// updateState will "forgot" old failures if they exist.
func (cb *CircuitBreaker) updateState() {
	prev := time.Now().Add(-config.CircuitBreakerRefreshTime)

	for ; cb.lastFailureTime.Before(prev); {
		if cb.failuresTotal == 0 {
			break
		}

		cb.lastFailureTime.Add(config.CircuitBreakerRefreshTime)
		cb.failuresTotal--
	}

	now := time.Now()
	cb.lastFailureTime = &now
}
