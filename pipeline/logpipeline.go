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
	return &LogPipeline{reader: reader, lineChannel: make(parser.LogLineChannel),
		dateChannel: make(parser.LogLineChannel), columnChannel: make(parser.LogLineChannel)}
}

func (p *LogPipeline) GetOutput() parser.LogLineChannel {
	return p.columnChannel
}

func (p *LogPipeline) Start() {
	go parser.ParseColumns(p.columnChannel, p.dateChannel)
	go parser.ParseDates(p.dateChannel, p.lineChannel)
	parser.ParseLines(p.lineChannel, p.reader)
}

func (p *LogPipeline) Close() {
	close(p.lineChannel)
	close(p.dateChannel)
	close(p.columnChannel)
}
