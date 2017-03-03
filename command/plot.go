package command

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/subcommands"
	"github.com/kbence/logan/config"
	"github.com/kbence/logan/pipeline"
	"github.com/kbence/logan/types"
	"github.com/kbence/logan/utils"
)

type plotCmd struct {
	config       *config.Configuration
	timeInterval string
	fields       string
	mode         string
	autoUpdate   bool
}

// PlotCommand creates a new plotCmd instance
func PlotCommand(config *config.Configuration) subcommands.Command {
	return &plotCmd{config: config}
}

func (c *plotCmd) Name() string {
	return "plot"
}

func (c *plotCmd) Synopsis() string {
	return "Plots a time-based line chart from the number of log lines"
}

func (c *plotCmd) Usage() string {
	return "plot <category>:\n" +
		"    draws a line chart from the number of log lines\n"
}

func (c *plotCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.timeInterval, "t", "-1h", "Example: -1h5m+5m")
	f.StringVar(&c.fields, "f", "", "Example: 1,2,3")
	f.StringVar(&c.mode, "m", "braille",
		fmt.Sprintf("One of the following modes: %s.", strings.Join(types.CharacterSets.GetNames(), ", ")))
	f.BoolVar(&c.autoUpdate, "u", true, "Auto-update chart during log parsing")
}

func (c *plotCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()

	if len(args) < 1 {
		log.Fatal("You have to pass a log source to this command!")
	}

	interval := utils.ParseTimeInterval(c.timeInterval, time.Now())
	width, height := utils.GetTerminalDimensions()

	p := pipeline.NewPipelineBuilder(pipeline.PipelineSettings{
		Category: args[0],
		Interval: interval,
		Filters:  args[1:],
		Fields:   utils.ParseIntervals(c.fields),
		Config:   c.config,
		Output:   pipeline.OutputTypeLineChart,
		OutputSettings: pipeline.LineChartSettings{
			Mode:            c.mode,
			Width:           width,
			Height:          height - 1,
			Interval:        interval,
			FrequentUpdates: c.autoUpdate}})
	p.Execute()

	return subcommands.ExitSuccess
}
