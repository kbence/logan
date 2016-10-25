package parser

import (
	"bufio"
	"io"
	"log"
)

func ParseLines(reader io.Reader, output chan string) {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')

		if line != "" {
			if line[len(line)-1] == '\n' {
				line = line[0 : len(line)-1]
			}

			output <- line
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	}
}
