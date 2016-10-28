package pipeline

import (
	"github.com/kbence/logan/filter"
	"github.com/kbence/logan/parser"
)

type FilterPipeline struct {
	filters      []filter.Filter
	channels     []parser.LogLineChannel
	inputChannel parser.LogLineChannel
}

func NewFilterPipeline(input parser.LogLineChannel, filters []filter.Filter) *FilterPipeline {
	return &FilterPipeline{filters: filters, inputChannel: input}
}

func filterFunc(output parser.LogLineChannel, input parser.LogLineChannel, f filter.Filter) {
	for {
		line, more := <-input

		if !more {
			break
		}

		if f.Match(line) {
			output <- line
		}
	}
}

func (p *FilterPipeline) Start() parser.LogLineChannel {
	inputChannel := p.inputChannel

	for _, f := range p.filters {
		outputChannel := make(parser.LogLineChannel)
		go filterFunc(outputChannel, inputChannel, f)
		inputChannel = outputChannel
	}

	return inputChannel
}
