package pipeline

import (
	"fmt"

	"github.com/kbence/logan/types"
)

type LineChartSettings struct {
	Mode     string
	Width    int
	Height   int
	Interval *types.TimeInterval
}

type LineChartPipeline struct {
	input         types.LogLineChannel
	settings      LineChartSettings
	chartSettings *types.ChartSettings
	sampler       *types.TimelineSampler
}

func NewLineChartPipeline(input types.LogLineChannel, settings LineChartSettings) *LineChartPipeline {
	chartSettings := &types.ChartSettings{
		Mode:        settings.Mode,
		Width:       settings.Width - 1,
		Height:      settings.Height - 1,
		Border:      true,
		YAxisLabels: true,
		XAxisLabels: true,
		Interval:    settings.Interval,
	}

	sampler := types.NewTimelineSampler(settings.Interval, chartSettings.SamplerSize())

	return &LineChartPipeline{
		input:         input,
		settings:      settings,
		chartSettings: chartSettings,
		sampler:       sampler,
	}
}

func (p *LineChartPipeline) render() {
	chartRenderer := types.NewChartRenderer(p.chartSettings)
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
