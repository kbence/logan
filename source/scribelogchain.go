package source

import (
	"fmt"
	"io"
	"log"
	"os"
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
