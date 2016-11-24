package parser

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/kbence/logan/types"
)

type preprocessorFunc func(string) string

type dateParser struct {
	reg          *regexp.Regexp
	layout       string
	preprocessor preprocessorFunc
}

var dateParsers = []dateParser{
	dateParser{
		reg:    regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}[.,]\\d{3})"),
		layout: "2006-01-02 15:04:05.000",
		preprocessor: func(input string) string {
			return strings.Replace(input, ",", ".", -1)
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}[.,]\\d{3}[+-]\\d{4})"),
		layout: "2006-01-02T15:04:05.000-0700",
		preprocessor: func(input string) string {
			return strings.Replace(input, ",", ".", -1)
		}}}

var location *time.Location

func getLocation() *time.Location {
	if location == nil {
		var err error

		location, err = time.LoadLocation("Local")

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

			return &parsedDate
		}
	}

	return nil
}

func ParseDates(output types.LogLineChannel, input types.LogLineChannel) {
	for {
		line, more := <-input

		if !more {
			break
		}

		if date := ParseDate(line.Line); date != nil {
			line.Date = *date
		}

		output <- line
	}

	close(output)
}
