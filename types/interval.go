package types

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TimeInterval struct {
	StartTime time.Time
	EndTime   time.Time
}

func NewTimeInterval(startTime, endTime time.Time) *TimeInterval {
	return &TimeInterval{StartTime: startTime, EndTime: endTime}
}

type IntInterval struct {
	Start int
	End   int
}

func NewIntInterval(start, end int) *IntInterval {
	return &IntInterval{Start: start, End: end}
}

var singleIntIntervalMatcher = regexp.MustCompile("^[0-9]+$")
var fromToIntIntervalMatcher = regexp.MustCompile("^([0-9]+)-([0-9]+)$")

func ParseIntInterval(interval string) *IntInterval {
	if singleIntIntervalMatcher.MatchString(interval) {
		value, _ := strconv.ParseInt(interval, 10, 31)
		return NewIntInterval(int(value), int(value))
	}

	if matches := fromToIntIntervalMatcher.FindStringSubmatch(interval); matches != nil {
		start, _ := strconv.ParseInt(matches[1], 10, 31)
		end, _ := strconv.ParseInt(matches[2], 10, 31)
		return NewIntInterval(int(start), int(end))
	}

	log.Panicf("Interval cannot be parsed: %s", interval)

	// Unreachable code
	return NewIntInterval(-1, -1)
}

func ParseAllIntIntervals(intervals string) []*IntInterval {
	results := []*IntInterval{}

	for _, interval := range strings.Split(intervals, ",") {
		results = append(results, ParseIntInterval(interval))
	}

	return results
}

func (i *IntInterval) String() string {
	return fmt.Sprintf("[Interval %d-%d]", i.Start, i.End)
}
