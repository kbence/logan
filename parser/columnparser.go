package parser

import "bytes"

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func ParseColumns(output LogLineChannel, input LogLineChannel) {
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
					line.Columns = append(line.Columns, columnBuffer.String())
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
			line.Columns = append(line.Columns, columnBuffer.String())
		}

		output <- line
	}
}
