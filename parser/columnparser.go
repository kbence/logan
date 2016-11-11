package parser

import (
	"bytes"

	"github.com/kbence/logan/types"
)

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func ParseColumns(output types.LogLineChannel, input types.LogLineChannel) {
	for {
		line, more := <-input

		if !more {
			close(output)
			break
		}

		lineLen := len(line.Line)
		ws := true
		quoted := false
		columnBuffer := bytes.Buffer{}
		currentColumn := 1
		line.Columns = map[int]string{}

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
					line.Columns[currentColumn] = columnBuffer.String()
					currentColumn++
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
			line.Columns[currentColumn] = columnBuffer.String()
		}

		output <- line
	}
}
