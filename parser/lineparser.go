package parser

import (
	"bufio"
	"io"
	"log"

	"github.com/kbence/logan/types"
)

func ParseLines(output types.LogLineChannel, reader io.Reader) {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')

		if line != "" {
			if line[len(line)-1] == '\n' {
				line = line[0 : len(line)-1]
			}

			output <- &types.LogLine{Line: line}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	}

	close(output)
}
