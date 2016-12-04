package filter

import "github.com/kbence/logan/types"

type TimeFilter struct {
	Interval *types.TimeInterval
}

func NewTimeFilter(interval *types.TimeInterval) *TimeFilter {
	return &TimeFilter{Interval: interval}
}

func (f *TimeFilter) Match(line *types.LogLine) bool {
	return f.Interval.Contains(line.Date)
}
