package command

import (
	"flag"
	"log"
	"time"

	"github.com/google/subcommands"
	"github.com/kbence/logan/config"
	"github.com/kbence/logan/pipeline"
	"github.com/kbence/logan/utils"
	"golang.org/x/net/context"
)

type inspectCmd struct {
	config       *config.Configuration
	timeInterval string
}

// InspectCommand creates a new inspectCmd instance
func InspectCommand(config *config.Configuration) subcommands.Command {
	return &inspectCmd{config: config}
}

func (c *inspectCmd) Name() string {
	return "inspect"
}

func (c *inspectCmd) Synopsis() string {
	return "Inspects the first 5 lines of log"
}

func (c *inspectCmd) Usage() string {
	return "inspect <category>:\n" +
		"    helper command to show field values in log files\n"
}

func (c *inspectCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.timeInterval, "t", "-1h", "Example: -1h5m+5m")
}

func (c *inspectCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()

	if len(args) < 1 {
		log.Fatal("You have to pass a log source to this command!")
	}

	p := pipeline.NewPipelineBuilder(pipeline.PipelineSettings{
		Category: args[0],
		Interval: utils.ParseTimeInterval(c.timeInterval, time.Now()),
		Filters:  args[1:],
		Fields:   utils.ParseIntervals(""),
		Config:   c.config,
		Output:   pipeline.OutputTypeInspector})
	p.Execute()

	return subcommands.ExitSuccess
}
