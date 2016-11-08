package pipeline

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/kbence/logan/config"
	"github.com/kbence/logan/filter"
	"github.com/kbence/logan/source"
	"github.com/kbence/logan/types"
)

type PipelineSettings struct {
	Category string
	Interval *types.TimeInterval
	Filters  []string
	Fields   []*types.IntInterval
	Config   *config.Configuration
}

type PipelineBuilder struct {
	settings PipelineSettings
}

func NewPipelineBuilder(settings PipelineSettings) *PipelineBuilder {
	fmt.Println(settings.Fields)
	return &PipelineBuilder{settings: settings}
}

func (p *PipelineBuilder) Execute() {
	categoryParts := strings.Split(p.settings.Category, "/")

	if len(categoryParts) != 2 {
		log.Fatal("Log source must be in the following format: source/category!")
	}

	src := categoryParts[0]
	category := categoryParts[1]

	logSource := source.GetLogSource(p.settings.Config, src)
	chain := logSource.GetChain(category)

	if chain == nil {
		log.Fatalf("Category '%s' not found!", category)
	}

	filters := []filter.Filter{
		filter.NewTimeFilter(p.settings.Interval.StartTime, p.settings.Interval.EndTime),
	}

	for _, filterString := range p.settings.Filters {
		columnFilter := filter.NewColumnFilter(filterString)
		filters = append(filters, columnFilter)
	}

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()

	logPipeline := NewLogPipeline(reader)

	filterPipeline := NewFilterPipeline(logPipeline.Start(), filters)
	transformPipeline := NewTransformPipeline(filterPipeline.Start(), p.settings.Fields)
	outputPipeline := NewOutputPipeline(transformPipeline.Start())
	outputPipeline.Start()

	io.Copy(writer, chain.Between(p.settings.Interval.StartTime, p.settings.Interval.EndTime))
}
