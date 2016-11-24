package pipeline

import (
	"log"
	"strings"

	"github.com/kbence/logan/config"
	"github.com/kbence/logan/filter"
	"github.com/kbence/logan/source"
	"github.com/kbence/logan/types"
)

type OutputType int

const (
	OutputTypeLogLines OutputType = iota
	OutputTypeUniqueLines
	OutputTypeInspector
	OutputTypeLineChart
)

type PipelineSettings struct {
	Category       string
	Interval       *types.TimeInterval
	Filters        []string
	Fields         []*types.IntInterval
	Config         *config.Configuration
	Output         OutputType
	OutputSettings interface{}
}

type PipelineBuilder struct {
	settings PipelineSettings
}

func NewPipelineBuilder(settings PipelineSettings) *PipelineBuilder {
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

	logReader := chain.Between(p.settings.Interval.StartTime, p.settings.Interval.EndTime)
	logPipeline := NewLogPipeline(NewTimeAwareBufferedReader(logReader, p.settings.Interval))

	filterPipeline := NewFilterPipeline(logPipeline.Start(), filters)
	transformPipeline := NewTransformPipeline(filterPipeline.Start(), p.settings.Fields)

	var outputPipeline OutputPipeline

	switch p.settings.Output {
	case OutputTypeLogLines:
		outputPipeline = NewLogPrinterPipeline(transformPipeline.Start())
		break

	case OutputTypeUniqueLines:
		outputPipeline = NewUniquePipeline(transformPipeline.Start())
		break

	case OutputTypeInspector:
		outputPipeline = NewInspectPipeline(transformPipeline.Start())
		break

	case OutputTypeLineChart:
		outputPipeline = NewLineChartPipeline(transformPipeline.Start(),
			p.settings.OutputSettings.(LineChartSettings))
		break
	}

	<-outputPipeline.Start()
}
