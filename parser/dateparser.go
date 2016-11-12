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

func ParseDates(output types.LogLineChannel, input types.LogLineChannel) {
	location, err := time.LoadLocation("Local")

	if err != nil {
		log.Panicf("ERROR loading local timezone: %s", err)
	}

	for {
		line, more := <-input

		if !more {
			break
		}

		for _, parser := range dateParsers {
			date := parser.reg.FindString(line.Line)

			if date != "" {
				parsedDate, err := time.ParseInLocation(parser.layout, parser.preprocessor(date), location)

				if err != nil {
					continue
				}

				line.Date = parsedDate
				break
			}
		}

		output <- line
	}

	close(output)
}
