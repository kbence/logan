package parser

import "time"

type LogLine struct {
	Line    string
	Date    time.Time
	Columns map[int]string
}

type LogLineChannel chan *LogLine
