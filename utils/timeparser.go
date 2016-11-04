package utils

import (
	"bytes"
	"log"
	"time"

	"github.com/kbence/logan/types"
)

func breakTimeIntervalSpec(timeSpec string) []string {
	components := []*bytes.Buffer{bytes.NewBufferString(""), bytes.NewBufferString(""),
		bytes.NewBufferString("")}
	currentComponent := 0

	for _, ch := range []byte(timeSpec) {
		switch ch {
		case '+':
			currentComponent = 2
			break

		case '-':
			currentComponent = 1
			break

		default:
			components[currentComponent].Write([]byte{ch})
		}
	}

	return []string{components[0].String(), components[1].String(), components[2].String()}
}

type timeComponentArray [7]int
type parserFunc func(spec string) (timeComponentArray, error)

var parserFunctions = []parserFunc{
	func(spec string) (timeComponentArray, error) {
		parsed, err := time.Parse("15:04", spec)
		if err != nil {
			return timeComponentArray{}, err
		}
		return timeComponentArray{-1, -1, -1, parsed.Hour(), parsed.Minute(), 0, 0}, nil
	},
	func(spec string) (timeComponentArray, error) {
		parsed, err := time.Parse("15:04:05", spec)
		if err != nil {
			return timeComponentArray{}, err
		}
		return timeComponentArray{-1, -1, -1, parsed.Hour(), parsed.Minute(), parsed.Second(), 0}, nil
	}}

func parseTimeReference(ref string, now time.Time) time.Time {
	timeComponents := []int{now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(),
		now.Second(), now.Nanosecond()}

	for _, parse := range parserFunctions {
		parsedTime, err := parse(ref)

		if err != nil {
			continue
		}

		for i, comp := range parsedTime {
			if comp >= 0 {
				timeComponents[i] = comp
			}
		}
	}

	// Correct time from last 24h
	result := time.Date(timeComponents[0], time.Month(timeComponents[1]), timeComponents[2],
		timeComponents[3], timeComponents[4], timeComponents[5], timeComponents[6],
		now.Location())

	if result.After(now) {
		result = result.Add(-24 * time.Hour)
	}

	return result
}

func ParseTimeInterval(timeIntVal string, now time.Time) *types.TimeInterval {
	var timeRef, start, end time.Time
	components := breakTimeIntervalSpec(timeIntVal)

	if components[0] != "" {
		timeRef = parseTimeReference(components[0], now)
	} else {
		timeRef = now
	}

	if components[1] != "" {
		startDur, err := time.ParseDuration(components[1])
		if err != nil {
			log.Panicf("ERROR parsing time interval \"%s\": %s", components[1], err)
		}
		start = timeRef.Add(-startDur)
	} else {
		start = timeRef
	}

	if components[2] != "" {
		endDur, err := time.ParseDuration(components[2])
		if err != nil {
			log.Panicf("ERROR parsing time interval \"%s\": %s", components[2], err)
		}

		end = start.Add(endDur)
	} else {
		end = now
	}

	return types.NewTimeInterval(start, end)
}
