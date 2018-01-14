package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kbence/logan/types"
)

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isJSONLine(line string) bool {
	trimmedLine := strings.Trim(line, "\t ")
	llen := len(trimmedLine)

	return llen >= 2 && trimmedLine[0] == '{' && trimmedLine[llen-1:] == "}"
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

func parseLogLine(line *types.LogLine) {
	lineLen := len(line.Line)
	ws := true
	quoted := false
	quoteType := 0
	columnBuffer := bytes.Buffer{}
	currentColumn := 1
	line.Columns = map[string]string{}

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
				line.Columns[fmt.Sprint(currentColumn)] = columnBuffer.String()
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
		line.Columns[fmt.Sprint(currentColumn)] = columnBuffer.String()
	}
}

func parseJSONLine(line *types.LogLine) {
	var object interface{}
	json.Unmarshal([]byte(line.String()), &object)

	line.Columns = types.ColumnList{}
	for key, value := range object.(map[string]interface{}) {
		line.Columns[key] = value.(string)
	}
}

func ParseColumns(output types.LogLineChannel, input types.LogLineChannel) {
	for {
		line, more := <-input

		if !more {
			break
		}

		if isJSONLine(line.Line) {
			parseJSONLine(line)
		} else {
			parseLogLine(line)
		}

		output <- line
	}

	close(output)
}
