package types

import "time"

type TimeInterval struct {
	StartTime time.Time
	EndTime   time.Time
}

func NewTimeInterval(startTime, endTime time.Time) *TimeInterval {
	return &TimeInterval{StartTime: startTime, EndTime: endTime}
}
