package parser

import (
	"log"
	"regexp"
	"time"
)

type dateParser struct {
	reg    *regexp.Regexp
	layout string
}

var dateParsers = []dateParser{
	dateParser{
		reg:    regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}.\\d{3})"),
		layout: "2006-01-02 03:04:05.000"}}

func ParseDates(input LogLineChannel, output LogLineChannel) {
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
				parsedDate, err := time.ParseInLocation(parser.layout, date, location)

				if err != nil {
					continue
				}

				line.Date = parsedDate
				break
			}
		}

		output <- line

	}
}
