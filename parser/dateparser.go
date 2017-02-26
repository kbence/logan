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

var multiSpaceMatcher = regexp.MustCompile("\\s+")
var commaMatcher = regexp.MustCompile(",")

var dateParsers = []dateParser{
	dateParser{
		reg:    regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}[.,]\\d{6}[+-] ([+-]\\d{4}|[A-Z]{3}))"),
		layout: "2006-01-02T15:04:05.000000 -0700",
		preprocessor: func(input string) string {
			return strings.Replace(input, ",", ".", -1)
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}[.,]\\d{3}[+-]\\d{4})"),
		layout: "2006-01-02T15:04:05.000-0700",
		preprocessor: func(input string) string {
			return strings.Replace(input, ",", ".", -1)
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}[.,]\\d{3})"),
		layout: "2006-01-02 15:04:05.000",
		preprocessor: func(input string) string {
			return strings.Replace(input, ",", ".", -1)
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\w{3} ( ?\\d|\\d{2}) \\d{2}:\\d{2}:\\d{2} \\d{4})"),
		layout: "Jan 2 15:04:05 2006",
		preprocessor: func(input string) string {
			return multiSpaceMatcher.ReplaceAllString(input, " ")
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\w{3} ( ?\\d|\\d{2}) \\d{2}:\\d{2}:\\d{2})"),
		layout: "Jan 2 15:04:05",
		preprocessor: func(input string) string {
			return multiSpaceMatcher.ReplaceAllString(input, " ")
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\w{3},? \\d{2} \\w{3} \\d{2}:\\d{2}:\\d{2}.\\d{3} [+-]\\d{4})"),
		layout: "Mon 02 Jan 15:04:05.000 -0700",
		preprocessor: func(input string) string {
			return commaMatcher.ReplaceAllString(input, "")
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\w{3},? \\d{2} \\w{3} \\d{2}:\\d{2}:\\d{2}.\\d{3} [A-Z]{3})"),
		layout: "Mon 02 Jan 15:04:05.000 MST",
		preprocessor: func(input string) string {
			return commaMatcher.ReplaceAllString(input, "")
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\w{3},? \\d{2} \\w{3} \\d{2}:\\d{2}:\\d{2} [+-]\\d{4})"),
		layout: "Mon 02 Jan 15:04:05 -0700",
		preprocessor: func(input string) string {
			return commaMatcher.ReplaceAllString(input, "")
		}},
	dateParser{
		reg:    regexp.MustCompile("^(\\w{3},? \\d{2} \\w{3} \\d{2}:\\d{2}:\\d{2} [A-Z]{3})"),
		layout: "Mon 02 Jan 15:04:05 MST",
		preprocessor: func(input string) string {
			return commaMatcher.ReplaceAllString(input, "")
		}}}

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
