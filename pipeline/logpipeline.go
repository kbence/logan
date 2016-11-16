package pipeline

import (
	"io"

	"github.com/kbence/logan/parser"
	"github.com/kbence/logan/types"
)

type LogPipeline struct {
	reader        io.Reader
	lineChannel   types.LogLineChannel
	dateChannel   types.LogLineChannel
	columnChannel types.LogLineChannel
}

func NewLogPipeline(reader io.Reader) *LogPipeline {
	return &LogPipeline{reader: reader}
}

func (p *LogPipeline) GetOutput() types.LogLineChannel {
	return p.columnChannel
}

func (p *LogPipeline) Start() types.LogLineChannel {
	p.lineChannel = types.NewLogLineChannel()
	p.dateChannel = types.NewLogLineChannel()
	p.columnChannel = types.NewLogLineChannel()

	go parser.ParseColumns(p.columnChannel, p.dateChannel)
	go parser.ParseDates(p.dateChannel, p.lineChannel)
	go parser.ParseLines(p.lineChannel, p.reader)

	return p.columnChannel
}
