package pipeline

import (
	"fmt"

	"github.com/kbence/logan/types"
)

type UniquePipeline struct {
	input   types.LogLineChannel
	counter *types.LogLineCounter
}

func NewUniquePipeline(input types.LogLineChannel) *UniquePipeline {
	return &UniquePipeline{input: input, counter: types.NewLogLineCounter()}
}

func (p *UniquePipeline) storeUniqueLine(line *types.LogLine) {
	p.counter.Add(line)
}

func (p *UniquePipeline) printUniqueLines() {
	for _, line := range p.counter.UniqueLines(true) {
		fmt.Printf("%9d | ", line.Count)
		printColumnsInOrder(line.Line.Columns)
	}
}

func (p *UniquePipeline) Start() chan bool {
	exitChannel := make(chan bool)

	go func() {
		for {
			line, more := <-p.input

			if !more {
				break
			}

			p.storeUniqueLine(line)
		}

		p.printUniqueLines()
		exitChannel <- true
	}()

	return exitChannel
}
