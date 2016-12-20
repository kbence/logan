package command

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/google/subcommands"
	"github.com/kbence/logan/config"
	"github.com/kbence/logan/source"
)

type listCmd struct {
	config *config.Configuration
}

// ListCommand creates a new listCmd instance
func ListCommand(config *config.Configuration) subcommands.Command {
	return &listCmd{config: config}
}

func (c *listCmd) Name() string {
	return "list"
}

func (c *listCmd) Synopsis() string {
	return "Lists available log categories"
}

func (c *listCmd) Usage() string {
	return "list:\n" +
		"    print available log categories"
}

func (c *listCmd) SetFlags(f *flag.FlagSet) {}

func (c *listCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	allSources := []string{}

	for name, source := range source.GetLogSources(c.config) {
		for _, category := range source.GetCategories() {
			allSources = append(allSources, fmt.Sprintf("%s/%s", name, category))
		}
	}

	sort.Strings(allSources)
	fmt.Println(strings.Join(allSources, "\n"))

	return subcommands.ExitSuccess
}
