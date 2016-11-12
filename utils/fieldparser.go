package utils

import (
	"math"

	"github.com/kbence/logan/types"
)

// ParseIntervals returns parsed IntIntervals or 1-MaxInt32 for empty string
func ParseIntervals(fields string) []*types.IntInterval {
	if len(fields) > 0 {
		return types.ParseAllIntIntervals(fields)
	}

	return []*types.IntInterval{types.NewIntInterval(1, math.MaxInt32)}
}
