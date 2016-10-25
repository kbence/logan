package parser

import (
	"bytes"
	"time"
)

type ColumnsWithDate struct {
	Date    time.Time
	Columns []string
	Line    string
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func ParseColumns(input chan *LineWithDate, output chan *ColumnsWithDate) {
	for {
		line, more := <-input

		if !more {
			break
		}

		result := ColumnsWithDate{Date: line.Date, Columns: []string{}, Line: line.Line}

		lineLen := len(line.Line)
		ws := true
		quoted := false
		columnBuffer := bytes.Buffer{}

		for n := 0; n < lineLen; n++ {
			char := line.Line[n]
			isws := isWhitespace(char)

			if ws {
				if !isws {
					if char != '"' {
						columnBuffer.WriteByte(char)
					}
					ws = false
				}
			} else {
				if isws && !quoted {
					result.Columns = append(result.Columns, columnBuffer.String())
					columnBuffer = bytes.Buffer{}
					ws = true
				} else {
					if char != '"' {
						columnBuffer.WriteByte(char)
					}
				}
			}

			if char == '"' {
				quoted = !quoted
			}
		}

		if columnBuffer.Len() > 0 {
			result.Columns = append(result.Columns, columnBuffer.String())
		}

		output <- &result
	}
}
