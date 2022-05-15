// Package util stores some small and useful converters and formatters.
package util

import (
	"fmt"
	pb "github.com/elkopass/BITA/internal/proto"
	"strconv"
)

func QuotationToFloat(q pb.Quotation) float64 {
	if q.Nano < 0 {
		q.Nano = -q.Nano
	}
	p, _ := strconv.ParseFloat(fmt.Sprintf("%d.%d", q.Units, q.Nano), 64)

	return p
}

func MoneyValueToFloat(q pb.MoneyValue) float64 {
	if q.Nano < 0 {
		q.Nano = -q.Nano
	}
	p, _ := strconv.ParseFloat(fmt.Sprintf("%d.%d", q.Units, q.Nano), 64)

	return p
}
