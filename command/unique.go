package command

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/google/subcommands"
	"github.com/kbence/logan/config"
	"github.com/kbence/logan/pipeline"
	"github.com/kbence/logan/utils"
)

type uniqueCmd struct {
	config       *config.Configuration
	timeInterval string
	fields       string
	topLimit     int
}

// UniqCommand creates a new uniqueCmd instance
func UniqueCommand(config *config.Configuration) subcommands.Command {
	return &uniqueCmd{config: config}
}

func (c *uniqueCmd) Name() string {
	return "uniq"
}

func (c *uniqueCmd) Synopsis() string {
	return "Shows sum of unique lines"
}

func (c *uniqueCmd) Usage() string {
	return "list <category>:\n" +
		"    print log lines from the given category\n"
}

func (c *uniqueCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.timeInterval, "t", "-1h", "Example: -1h5m+5m")
	f.StringVar(&c.fields, "f", "", "Example: 1,2,3")
	f.IntVar(&c.topLimit, "top", 0, "Show only the top N results")
}

func (c *uniqueCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()

	if len(args) < 1 {
		log.Fatal("You have to pass a log source to this command!")
	}

	width, height := utils.GetTerminalDimensions()
	if c.topLimit > height-1 {
		c.topLimit = height - 1
	}

	p := pipeline.NewPipelineBuilder(pipeline.PipelineSettings{
		Category: args[0],
		Interval: utils.ParseTimeInterval(c.timeInterval, time.Now()),
		Filters:  args[1:],
		Fields:   utils.ParseIntervals(c.fields),
		Config:   c.config,
		Output:   pipeline.OutputTypeUniqueLines,
		OutputSettings: pipeline.UniqueSettings{
			TopLimit:      c.topLimit,
			TerminalWidth: width,
		}})
	p.Execute()

	return subcommands.ExitSuccess
}
