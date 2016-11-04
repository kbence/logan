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

type showCmd struct {
	config       *config.Configuration
	timeInterval string
}

// ShowCommand creates a new showCmd instance
func ShowCommand(config *config.Configuration) subcommands.Command {
	return &showCmd{config: config}
}

func (c *showCmd) Name() string {
	return "show"
}

func (c *showCmd) Synopsis() string {
	return "Shows log lines from the given category"
}

func (c *showCmd) Usage() string {
	return "list <category>:\n" +
		"    print log lines from the given category\n"
}

func (c *showCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.timeInterval, "t", "-1h", "Example: -1h5m+5m")
}

func (c *showCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()

	if len(args) < 1 {
		log.Fatal("You have to pass a log source to this command!")
	}

	p := pipeline.NewPipelineBuilder(pipeline.PipelineSettings{
		Category: args[0],
		Interval: utils.ParseTimeInterval(c.timeInterval, time.Now()),
		Filters:  args[1:],
		Config:   c.config})
	p.Execute()

	return subcommands.ExitSuccess
}
