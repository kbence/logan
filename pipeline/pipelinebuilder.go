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

func (p *PipelineBuilder) getChain() (chain source.LogChain) {
	var logSource source.LogSource

	switch strings.Count(p.settings.Category, "/") {
	case 0:
		sources := source.GetSourcesForCategory(p.settings.Config, p.settings.Category)

		if len(sources) > 1 {
			log.Fatalf("Ambiguous category name: '%s'! Please specify source name!", p.settings.Category)
		} else if len(sources) == 0 {
			log.Fatalf("Catergory '%s' not found!", p.settings.Category)
		}

		chain = sources[0].GetChain(p.settings.Category)
		break

	case 1:
		categoryParts := strings.Split(p.settings.Category, "/")
		src := categoryParts[0]
		category := categoryParts[1]

		logSource = source.GetLogSource(p.settings.Config, src)
		chain = logSource.GetChain(category)
		break

	default:
		log.Fatal("Log source must be in the following format: [source/]category!")
	}

	if chain == nil {
		log.Fatalf("Category '%s' not found!", p.settings.Category)
	}

	return chain
}

func (p *PipelineBuilder) Execute() {
	chain := p.getChain()

	filters := []filter.Filter{
		filter.NewTimeFilter(p.settings.Interval),
	}

	for _, filterString := range p.settings.Filters {
		columnFilter := filter.NewColumnFilter(filterString)
		filters = append(filters, columnFilter)
	}

	logReader := chain.Between(p.settings.Interval)
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
