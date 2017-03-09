package command

import (
	"log"
	"time"

	"github.com/kbence/logan/config"
	"github.com/kbence/logan/pipeline"
	"github.com/kbence/logan/utils"
	"github.com/spf13/cobra"
)

// NewUniqCommand returns the command responsible for creating the
// uniq output
func NewUniqCommand(cfg *config.Configuration) *cobra.Command {
	var timeInterval string
	var fields string
	var topLimit int

	uniqCommand := &cobra.Command{
		Use:   "uniq",
		Short: "Shows sum of unique lines",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("You have to pass a log source to this command!")
			}

			width, height := utils.GetTerminalDimensions()
			if topLimit > height-1 {
				topLimit = height - 1
			}

			p := pipeline.NewPipelineBuilder(pipeline.PipelineSettings{
				Category: args[0],
				Interval: utils.ParseTimeInterval(timeInterval, time.Now()),
				Filters:  args[1:],
				Fields:   utils.ParseIntervals(fields),
				Config:   cfg,
				Output:   pipeline.OutputTypeUniqueLines,
				OutputSettings: pipeline.UniqueSettings{
					TopLimit:      topLimit,
					TerminalWidth: width,
				}})
			p.Execute()
		},
	}

	uniqCommand.Flags().StringVarP(&timeInterval, "time", "t", "-1h", "Example: -1h5m+5m")
	uniqCommand.Flags().StringVarP(&fields, "fields", "f", "", "Example: 1,2,3")
	uniqCommand.Flags().IntVarP(&topLimit, "top", "T", 0, "Show only the top N results")

	return uniqCommand
}
