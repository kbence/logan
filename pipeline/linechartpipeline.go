package pipeline

import (
	"fmt"
	"os"
	"time"

	"github.com/kbence/logan/types"
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
	input          types.LogLineChannel
	settings       LineChartSettings
	chartSettings  *types.ChartSettings
	sampler        *types.TimelineSampler
	linesReceived  uint64
	nextUpdateLine uint64
	lastUpdateLine uint64
	lastUpdateTime *time.Time
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

func (p *LineChartPipeline) updateIfNeeded() {
	if !p.settings.FrequentUpdates {
		return
	}

	if p.lastUpdateTime == nil || p.linesReceived >= p.nextUpdateLine {
		p.render()
		fmt.Print(terminfo.GoUpBy(p.settings.Height - 1))
		os.Stdout.Sync()

		currentTime := time.Now()

		if p.lastUpdateTime == nil {
			p.nextUpdateLine = p.linesReceived + 10
		} else {
			next := 1000000000 * (p.linesReceived - p.lastUpdateLine) /
				uint64(currentTime.Sub(*p.lastUpdateTime))
			p.nextUpdateLine = p.linesReceived + next
		}

		p.lastUpdateTime = &currentTime
		p.lastUpdateLine = p.linesReceived
	}
}

func (p *LineChartPipeline) Start() chan bool {
	exitChannel := make(chan bool)

	go func() {
		p.linesReceived = 0

		for {
			line, more := <-p.input

			if !more {
				break
			}

			p.sampler.Inc(line.Date, 1)
			p.linesReceived++
			p.updateIfNeeded()
		}

		p.render()

		exitChannel <- true
	}()

	return exitChannel
}
