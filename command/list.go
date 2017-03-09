package command

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/kbence/logan/config"
	"github.com/kbence/logan/source"
	"github.com/spf13/cobra"
)

// NewListCommand creates a cobra.Command instance that implements the `list` command
func NewListCommand(cfg *config.Configuration) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists available log categories",
		Run: func(cmd *cobra.Command, args []string) {
			allSources := []string{}

			if len(args) > 1 {
				log.Fatal("List can only take 0 or 1 arguments!")
			}

			for name, source := range source.GetLogSources(cfg) {
				for _, category := range source.GetCategories() {
					allSources = append(allSources, fmt.Sprintf("%s/%s", name, category))
				}
			}

			if len(args) == 1 {
				filter := args[0]
				filteredSources := []string{}

				for _, source := range allSources {
					if strings.Contains(source, filter) {
						filteredSources = append(filteredSources, source)
					}
				}

				allSources = filteredSources
			}

			sort.Strings(allSources)
			fmt.Println(strings.Join(allSources, "\n"))
		},
	}
}
