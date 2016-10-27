package parser

import "time"

type LogLine struct {
	Line    string
	Date    time.Time
	Columns []string
}

type LogLineChannel chan *LogLine
