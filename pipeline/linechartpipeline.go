package pipeline

import (
	"fmt"
	"os"
	"sync"
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
	input         types.LogLineChannel
	settings      LineChartSettings
	chartSettings *types.ChartSettings
	sampler       *types.TimelineSampler
	renderLock    *sync.Mutex
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
		renderLock:    &sync.Mutex{},
	}
}

func (p *LineChartPipeline) render() {
	chartRenderer := types.NewChartRenderer(p.chartSettings)
	chartRenderer.AddDataLine(p.sampler.Samples)
	fmt.Println(chartRenderer.Render())
}

func (p *LineChartPipeline) updateLoop(exitChannel chan bool) {
	for {
		select {
		case <-time.After(time.Second):
			p.renderLock.Lock()
			p.render()
			fmt.Print(terminfo.GoUpBy(p.settings.Height - 1))
			os.Stdout.Sync()
			p.renderLock.Unlock()
			break

		case <-exitChannel:
			return
		}
	}
}

func (p *LineChartPipeline) Start() chan bool {
	exitChannel := make(chan bool)
	updateExitChannel := make(chan bool)

	if p.settings.FrequentUpdates {
		go p.updateLoop(updateExitChannel)
	}

	go func() {
		for {
			line, more := <-p.input

			if !more {
				break
			}

			p.sampler.Inc(line.Date, 1)
		}

		if p.settings.FrequentUpdates {
			// Exit update goproc _before_ we start rendering
			updateExitChannel <- true
		}

		p.renderLock.Lock()
		p.render()
		p.renderLock.Unlock()

		exitChannel <- true
	}()

	return exitChannel
}
