package pipeline

import (
	"fmt"
	"os"

	"github.com/kbence/logan/types"
	"github.com/kbence/logan/utils"
	"github.com/kbence/logan/utils/terminfo"
)

type LineChartSettings struct {
	Mode            string
	Width           int
	Height          int
	Interval        *types.TimeInterval
	FrequentUpdates bool
}

type LineChartPipeline struct {
	input         types.LogLineChannel
	settings      LineChartSettings
	chartSettings *types.ChartSettings
	sampler       *types.TimelineSampler
	timer         *utils.UpdateTimer
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
		timer:         utils.NewUpdateTimer(),
	}
}

func (p *LineChartPipeline) render() {
	chartRenderer := types.NewChartRenderer(p.chartSettings)
	chartRenderer.AddDataLine(p.sampler.Samples)
	fmt.Println(chartRenderer.Render())
}

func (p *LineChartPipeline) updateIfNeeded() {
	if !p.settings.FrequentUpdates {
		return
	}

	if p.timer.IsUpdateNeeded() {
		p.render()
		fmt.Print(terminfo.GoUpBy(p.settings.Height - 1))
		os.Stdout.Sync()
	}
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
			p.updateIfNeeded()
		}

		p.render()

		exitChannel <- true
	}()

	return exitChannel
}
