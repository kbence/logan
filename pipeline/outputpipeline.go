package pipeline

import (
	"fmt"

	"github.com/kbence/logan/parser"
)

type OutputPipeline struct {
	input parser.LogLineChannel
}

func NewOutputPipeline(input parser.LogLineChannel) *OutputPipeline {
	return &OutputPipeline{input: input}
}

func (p *OutputPipeline) Start() {
	go func() {
		for {
			line, more := <-p.input

			if !more {
				break
			}

			fmt.Println(line.Line)
		}
	}()
}
