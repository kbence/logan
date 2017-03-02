package parser

import (
	"regexp"
	"strings"
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
		reg:          regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}[.,]\\d{6}[+-] ([+-]\\d{4}|[A-Z]{3}))"),
		layout:       "2006-01-02T15:04:05.000000 -0700",
		preprocessor: replaceCommaToPoint},
	dateParser{
		reg:          regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}[.,]\\d{3}[+-]\\d{4})"),
		layout:       "2006-01-02T15:04:05.000-0700",
		preprocessor: replaceCommaToPoint},
	dateParser{
		reg:          regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}[.,]\\d{3})"),
		layout:       "2006-01-02 15:04:05.000",
		preprocessor: replaceCommaToPoint},
	dateParser{
		reg:          regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2})"),
		layout:       "2006-01-02 15:04:05",
		preprocessor: noop},
	dateParser{
		reg:          regexp.MustCompile("^(\\w{3} ( ?\\d|\\d{2}) \\d{2}:\\d{2}:\\d{2} \\d{4})"),
		layout:       "Jan 2 15:04:05 2006",
		preprocessor: normalizeSpaces},
	dateParser{
		reg:          regexp.MustCompile("^(\\w{3} ( ?\\d|\\d{2}) \\d{2}:\\d{2}:\\d{2})"),
		layout:       "Jan 2 15:04:05",
		preprocessor: normalizeSpaces},
	dateParser{
		reg:          regexp.MustCompile("^(\\w{3},? \\d{2} \\w{3} \\d{2}:\\d{2}:\\d{2}.\\d{3} [+-]\\d{4})"),
		layout:       "Mon 02 Jan 15:04:05.000 -0700",
		preprocessor: deleteComma},
	dateParser{
		reg:          regexp.MustCompile("^(\\w{3},? \\d{2} \\w{3} \\d{2}:\\d{2}:\\d{2}.\\d{3} [A-Z]{3})"),
		layout:       "Mon 02 Jan 15:04:05.000 MST",
		preprocessor: deleteComma},
	dateParser{
		reg:          regexp.MustCompile("^(\\w{3},? \\d{2} \\w{3} \\d{2}:\\d{2}:\\d{2} [+-]\\d{4})"),
		layout:       "Mon 02 Jan 15:04:05 -0700",
		preprocessor: deleteComma},
	dateParser{
		reg:          regexp.MustCompile("^(\\w{3},? \\d{2} \\w{3} \\d{2}:\\d{2}:\\d{2} [A-Z]{3})"),
		layout:       "Mon 02 Jan 15:04:05 MST",
		preprocessor: deleteComma},
	dateParser{
		reg:          regexp.MustCompile("^(\\w{3} \\w{3} ( \\d|\\d{2}) \\d{2}:\\d{2}:\\d{2}.\\d{3})"),
		layout:       "Mon Jan 2 15:04:05.000",
		preprocessor: normalizeSpaces}}

func noop(input string) string { return input }

func replaceCommaToPoint(input string) string {
	return strings.Replace(input, ",", ".", -1)
}

func deleteComma(input string) string {
	return commaMatcher.ReplaceAllString(input, "")
}

func normalizeSpaces(input string) string {
	return multiSpaceMatcher.ReplaceAllString(input, " ")
}
