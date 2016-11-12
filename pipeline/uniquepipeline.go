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

func drawPercentageBar(width int, percentage float64) string {
	barElements := []rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}
	barElementsLen := len(barElements)
	last := barElementsLen - 1
	bar := ""

	max := width * barElementsLen
	state := int(percentage * float64(max))

	for c := 0; c < width; c++ {
		var ch rune

		if c*barElementsLen < state {
			ch = barElements[last]
		} else if (c-1)*len(barElements) >= state {
			ch = barElements[0]
		} else {
			ch = barElements[(state-(c-1)*barElementsLen)-1]
		}

		bar = fmt.Sprintf("%s%c", bar, ch)
	}

	return bar
}

func (p *UniquePipeline) printUniqueLines() {
	max := p.counter.Max()

	for _, line := range p.counter.UniqueLines(true) {
		fmt.Printf("%9d|%s ", line.Count, drawPercentageBar(16, float64(line.Count)/float64(max)))
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
