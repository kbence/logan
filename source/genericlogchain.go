package source

import (
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
	"os"

	"github.com/kbence/logan/types"
)

type GenericLogChain struct {
	files []string
}

func NewGenericLogChain(files []string) LogChain {
	return &GenericLogChain{files: files}
}

func (c *GenericLogChain) Between(interval *types.TimeInterval) io.Reader {
	var readers []io.Reader

	for _, file := range c.files {
		var reader io.Reader
		reader, err := os.Open(file)

		if err != nil {
			log.Panicf("ERROR: %s\n", err)
		}

		if file[len(file)-3:] == ".gz" {
			reader, err = gzip.NewReader(reader)

			if err != nil {
				log.Panicf("ERROR: %s\n", err)
			}
		} else if file[len(file)-4:] == ".bz2" {
			reader = bzip2.NewReader(reader)
		}

		readers = append(readers, reader)
	}

	return io.MultiReader(readers...)
}
