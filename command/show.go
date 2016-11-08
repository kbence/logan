package command

import (
	"flag"
	"log"
	"math"
	"time"

	"github.com/google/subcommands"
	"github.com/kbence/logan/config"
	"github.com/kbence/logan/pipeline"
	"github.com/kbence/logan/types"
	"github.com/kbence/logan/utils"
	"golang.org/x/net/context"
)

type showCmd struct {
	config       *config.Configuration
	timeInterval string
	fields       string
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
	f.StringVar(&c.fields, "f", "", "Example: 1,2,3")
}

func parseIntervals(fields string) []*types.IntInterval {
	if len(fields) > 0 {
		return types.ParseAllIntIntervals(fields)
	}

	return []*types.IntInterval{types.NewIntInterval(1, math.MaxInt32)}
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
		Fields:   parseIntervals(c.fields),
		Config:   c.config})
	p.Execute()

	return subcommands.ExitSuccess
}
