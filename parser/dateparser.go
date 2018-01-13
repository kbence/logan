package parser

import (
	"log"
	"time"

	"github.com/kbence/logan/types"
)

var location *time.Location
var defaultLocationName = "Local"

func getLocation() *time.Location {
	if location == nil {
		var err error

		location, err = time.LoadLocation(defaultLocationName)

		if err != nil {
			log.Panicf("ERROR loading local timezone: %s", err)
		}
	}

	return location
}

func ParseDate(line string) *time.Time {
	for _, parser := range dateParsers {
		date := parser.reg.FindString(line)

		if date != "" {
			parsedDate, err := time.ParseInLocation(parser.layout, parser.preprocessor(date), getLocation())

			if err != nil {
				continue
			}

			if parsedDate.Year() == 0 {
				now := time.Now()
				parsedDate = time.Date(now.Year(), parsedDate.Month(), parsedDate.Day(),
					parsedDate.Hour(), parsedDate.Minute(), parsedDate.Second(),
					parsedDate.Nanosecond(), parsedDate.Location())

				if parsedDate.After(now) {
					parsedDate = time.Date(now.Year()-1, parsedDate.Month(), parsedDate.Day(),
						parsedDate.Hour(), parsedDate.Minute(), parsedDate.Second(),
						parsedDate.Nanosecond(), parsedDate.Location())
				}
			}
			return &parsedDate
		}
	}

	return nil
}

func ParseDates(output types.LogLineChannel, input types.LogLineChannel) {
	var lastDate *time.Time

	for {
		line, more := <-input

		if !more {
			break
		}

		if date := ParseDate(line.Line); date != nil {
			line.Date = *date
			lastDate = date
		} else if lastDate != nil {
			line.Date = *lastDate
		}

		output <- line
	}

	close(output)
}
