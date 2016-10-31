package pipeline

import (
	"io"

	"github.com/kbence/logan/parser"
)

type LogPipeline struct {
	reader        io.Reader
	lineChannel   parser.LogLineChannel
	dateChannel   parser.LogLineChannel
	columnChannel parser.LogLineChannel
}

func NewLogPipeline(reader io.Reader) *LogPipeline {
	return &LogPipeline{reader: reader}
}

func (p *LogPipeline) GetOutput() parser.LogLineChannel {
	return p.columnChannel
}

func (p *LogPipeline) Start() parser.LogLineChannel {
	p.lineChannel = make(parser.LogLineChannel)
	p.dateChannel = make(parser.LogLineChannel)
	p.columnChannel = make(parser.LogLineChannel)

	go parser.ParseColumns(p.columnChannel, p.dateChannel)
	go parser.ParseDates(p.dateChannel, p.lineChannel)
	go parser.ParseLines(p.lineChannel, p.reader)

	return p.columnChannel
}
