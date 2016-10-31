package parser

import (
	"bufio"
	"io"
	"log"
)

func ParseLines(output LogLineChannel, reader io.Reader) {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')

		if line != "" {
			if line[len(line)-1] == '\n' {
				line = line[0 : len(line)-1]
			}

			output <- &LogLine{Line: line}
		}

		if err == io.EOF {
			close(output)
			break
		}

		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	}
}
