package source

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type ScribeLogChain struct {
	directory string
	category  string
}

// NewScribeLogChain creates a new log chain for a given directory and
// scribe category
func NewScribeLogChain(directory string, category string) *ScribeLogChain {
	return &ScribeLogChain{directory: directory, category: category}
}

// Last returns a Reader to the latest log file in the chain
func (c *ScribeLogChain) Last() io.Reader {
	currentFile := fmt.Sprintf("%s/%s/%s_current", c.directory, c.category, c.category)

	if _, err := os.Stat(currentFile); os.IsNotExist(err) {
		log.Panicf("ERROR: %s\n", err)
	}

	realFile, err := os.Readlink(currentFile)

	if err != nil {
		log.Panicf("ERROR: %s\n", err)
	}

	file, err := os.Open(fmt.Sprintf("%s/%s/%s", c.directory, c.category, realFile))

	if err != nil {
		log.Panicf("ERROR: %s\n", err)
	}

	return file
}

func (c *ScribeLogChain) Between(start time.Time, end time.Time) io.Reader {
	readers := []io.Reader{}

	endHour := end.Hour()
	if end.Minute() > 0 || end.Second() > 0 {
		endHour++
	}

	endTime := time.Date(end.Year(), end.Month(), end.Day(), endHour, 0, 0, 0, end.Location())

	for curTime := start; curTime.Before(endTime); curTime = curTime.Add(time.Hour) {
		logtime := fmt.Sprintf("%s_000%s", curTime.Format("2006-01-02"), curTime.Format("15"))
		currentFile := fmt.Sprintf("%s/%s/%s-%s", c.directory, c.category, c.category, logtime)
		useGzip := false

		if _, err := os.Stat(currentFile); os.IsNotExist(err) {
			currentFile = fmt.Sprintf("%s.gz", currentFile)
			useGzip = true
		}

		var reader io.Reader
		reader, err := os.Open(currentFile)

		if err != nil {
			log.Printf("WARNING: %s\n", err)
			continue
		}

		if useGzip {
			reader, err = gzip.NewReader(reader)

			if err != nil {
				log.Panicf("ERROR: %s\n", err)
			}
		}

		readers = append(readers, reader)
		curTime.Add(time.Hour)
	}

	return io.MultiReader(readers...)
}
