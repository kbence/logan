package pipeline

import (
	"fmt"

	"github.com/kbence/logan/types"
)

type LineChartSettings struct {
	Width    int
	Height   int
	Interval *types.TimeInterval
}

type LineChartPipeline struct {
	input    types.LogLineChannel
	settings LineChartSettings
	sampler  *types.TimelineSampler
}

func NewLineChartPipeline(input types.LogLineChannel, settings LineChartSettings) *LineChartPipeline {
	sampler := types.NewTimelineSampler(settings.Interval, settings.Width*2)
	return &LineChartPipeline{input: input, settings: settings, sampler: sampler}
}

func (p *LineChartPipeline) render() {
	chartRenderer := types.NewChartRenderer(p.settings.Width, p.settings.Height)
	chartRenderer.AddDataLine(p.sampler.Samples)
	fmt.Println(chartRenderer.Render())
}

func (p *LineChartPipeline) Start() chan bool {
	exitChannel := make(chan bool)

	go func() {
		for {
			line, more := <-p.input

			if !more {
				break
			}

			p.sampler.Inc(line.Date, 1)
		}

		p.render()

		exitChannel <- true
	}()

	return exitChannel
}
