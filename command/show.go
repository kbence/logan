package command

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/google/subcommands"
	"github.com/kbence/logan/config"
	"github.com/kbence/logan/parser"
	"github.com/kbence/logan/source"
	"golang.org/x/net/context"
)

type showCmd struct {
	config *config.Configuration
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
		"    print log lines from the given category"
}

func (c *showCmd) SetFlags(f *flag.FlagSet) {
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

	columnChannel := make(chan *parser.ColumnsWithDate)
	dateChannel := make(chan *parser.LineWithDate)
	lineChannel := make(chan string)
	defer close(lineChannel)
	defer close(dateChannel)
	defer close(columnChannel)

	go func() {
		for {
			line := <-columnChannel
			fmt.Printf("%s\n", line.Line)
		}
	}()
	go parser.ParseColumns(dateChannel, columnChannel)
	go parser.ParseDates(lineChannel, dateChannel)
	parser.ParseLines(chain.Last(), lineChannel)

	return subcommands.ExitSuccess
}
