package command

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/subcommands"
	"github.com/kbence/logan/config"
	"github.com/kbence/logan/filter"
	"github.com/kbence/logan/parser"
	"github.com/kbence/logan/source"
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

func parseTimeInterval(timeIntVal string) (time.Time, time.Time) {
	var start, end time.Time
	parts := strings.Split(timeIntVal, "+")

	startDur, err := time.ParseDuration(parts[0])
	if err != nil {
		log.Panicf("ERROR parsing time interval \"%s\": %s", parts[0], err)
	}
	start = time.Now().Add(startDur)

	if len(parts) > 1 {
		endDur, err := time.ParseDuration(parts[1])
		if err != nil {
			log.Panicf("ERROR parsing time interval \"%s\": %s", parts[1], err)
		}

		end = start.Add(endDur)
	} else {
		end = time.Now()
	}

	return start, end
}

func (c *showCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()

	if len(args) != 1 {
		log.Fatal("You have to pass a log source to this command!")
	}

	categoryParts := strings.Split(args[0], "/")

	if len(categoryParts) != 2 {
		log.Fatal("Log source must be in the following format: source/category!")
	}

	src := categoryParts[0]
	category := categoryParts[1]

	logSource := source.GetLogSource(c.config, src)
	chain := logSource.GetChain(category)

	if chain == nil {
		log.Fatalf("Category '%s' not found!", category)
	}

	startTime, endTime := parseTimeInterval(c.timeInterval)

	timeFilter := filter.NewTimeFilter(startTime, endTime)
	columnChannel := make(parser.LogLineChannel)
	dateChannel := make(parser.LogLineChannel)
	lineChannel := make(parser.LogLineChannel)
	defer close(lineChannel)
	defer close(dateChannel)
	defer close(columnChannel)

	go func() {
		for {
			line, more := <-columnChannel

			if !more {
				return
			}

			if timeFilter.Match(line) {
				fmt.Printf("%s\n", line.Line)
			}
		}
	}()
	go parser.ParseColumns(dateChannel, columnChannel)
	go parser.ParseDates(lineChannel, dateChannel)
	parser.ParseLines(chain.Between(startTime, endTime), lineChannel)

	return subcommands.ExitSuccess
}
