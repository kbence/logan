package source

import (
	"io"
	"log"
	"os"

	"github.com/kbence/logan/types"
)

type StdinChain struct{}

func NewStdinChain(category string) *StdinChain {
	if category == "-" {
		return &StdinChain{}
	}

	log.Panic("ERROR: stdin category must be '-'!")
	return nil
}

func (c *StdinChain) Last() io.Reader {
	return os.Stdin
}

func (c *StdinChain) Between(interval *types.TimeInterval) io.Reader {
	return os.Stdin
}
