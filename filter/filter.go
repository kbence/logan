package filter

import "github.com/kbence/logan/parser"

type Filter interface {
	Match(line *parser.LogLine) bool
}
