package pipeline

import (
	"github.com/kbence/logan/filter"
	"github.com/kbence/logan/types"
)

type FilterPipeline struct {
	filters      []filter.Filter
	channels     []types.LogLineChannel
	inputChannel types.LogLineChannel
}

func NewFilterPipeline(input types.LogLineChannel, filters []filter.Filter) *FilterPipeline {
	return &FilterPipeline{filters: filters, inputChannel: input}
}

func filterFunc(output types.LogLineChannel, input types.LogLineChannel, f filter.Filter) {
	for {
		line, more := <-input

		if !more {
			close(output)
			break
		}

		if f.Match(line) {
			output <- line
		}
	}
}

func (p *FilterPipeline) Start() types.LogLineChannel {
	inputChannel := p.inputChannel

	for _, f := range p.filters {
		outputChannel := make(types.LogLineChannel)
		go filterFunc(outputChannel, inputChannel, f)
		inputChannel = outputChannel
	}

	return inputChannel
}
