package types

import "time"

type LogLine struct {
	Line    string
	Date    time.Time
	Columns ColumnList
}

type LogLineChannel chan *LogLine

func NewLogLineChannel() LogLineChannel {
	return make(LogLineChannel, 1024)
}

func (l *LogLine) String() string {
	return l.Line
}
