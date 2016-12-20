package parser

import (
	"bytes"

	"github.com/kbence/logan/types"
)

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

type separator struct {
	start, end byte
	keep       bool
}

var separators = []separator{
	separator{'"', '"', false},
	separator{'\'', '\'', false},
	separator{'[', ']', true},
	separator{'(', ')', true},
}

func isQuoteStart(char byte) bool {
	for _, sep := range separators {
		if sep.start == char {
			return true
		}
	}

	return false
}

func isQuoteEnd(char byte) bool {
	for _, sep := range separators {
		if sep.end == char {
			return true
		}
	}

	return false
}

func getQuoteType(startChar byte) int {
	for i, sep := range separators {
		if sep.start == startChar {
			return i
		}
	}

	return -1
}

func ParseColumns(output types.LogLineChannel, input types.LogLineChannel) {
	for {
		line, more := <-input

		if !more {
			break
		}

		lineLen := len(line.Line)
		ws := true
		quoted := false
		quoteType := 0
		columnBuffer := bytes.Buffer{}
		currentColumn := 1
		line.Columns = map[int]string{}

		for n := 0; n < lineLen; n++ {
			char := line.Line[n]
			isws := isWhitespace(char)

			if ws {
				if !isws {
					if !isQuoteStart(char) || separators[getQuoteType(char)].keep {
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
					if !isQuoteEnd(char) || separators[quoteType].keep {
						columnBuffer.WriteByte(char)
					}
				}
			}

			if quoted && char == separators[quoteType].end || !quoted && isQuoteStart(char) {
				quoted = !quoted

				if quoted {
					quoteType = getQuoteType(char)
				}
			}
		}

		if columnBuffer.Len() > 0 {
			line.Columns[currentColumn] = columnBuffer.String()
		}

		output <- line
	}

	close(output)
}
