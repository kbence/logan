package parser

import (
	"fmt"
	"regexp"
	"time"
)

type LineWithDate struct {
	Date time.Time
	Line string
}

type dateParser struct {
	reg    *regexp.Regexp
	layout string
}

var dateParsers = []dateParser{
	dateParser{
		reg:    regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}.\\d{3})"),
		layout: "2006-01-02 03:04:05.000"}}

func ParseDates(input chan string, output chan *LineWithDate) {
	for {
		line, more := <-input
		parsed := false

		if !more {
			break
		}

		for _, parser := range dateParsers {
			date := parser.reg.FindString(line)

			if date != "" {
				fmt.Println(date)
				parsedDate, err := time.Parse(parser.layout, date)

				if err != nil {
					continue
				}

				output <- &LineWithDate{Date: parsedDate, Line: line}
				parsed = true
				break
			}
		}

		if !parsed {
			output <- &LineWithDate{Date: time.Unix(0, 0), Line: line}
		}
	}
}
