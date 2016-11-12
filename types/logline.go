package types

import "time"

type LogLine struct {
	Line    string
	Date    time.Time
	Columns ColumnList
}

type LogLineChannel chan *LogLine
