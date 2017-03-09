package command

import (
	"log"
	"time"

	"github.com/kbence/logan/config"
	"github.com/kbence/logan/pipeline"
	"github.com/kbence/logan/utils"
	"github.com/spf13/cobra"
)

// NewShowCommand returns the command that is responsible for showing the
// (almost) raw lines of output
func NewShowCommand(cfg *config.Configuration) *cobra.Command {
	var timeInterval string
	var fields string

	showCommand := &cobra.Command{
		Use:   "show",
		Short: "Shows log lines from the given category",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("You have to pass a log source to this command!")
			}

			p := pipeline.NewPipelineBuilder(pipeline.PipelineSettings{
				Category: args[0],
				Interval: utils.ParseTimeInterval(timeInterval, time.Now()),
				Filters:  args[1:],
				Fields:   utils.ParseIntervals(fields),
				Config:   cfg,
				Output:   pipeline.OutputTypeLogLines})
			p.Execute()
		},
	}

	showCommand.Flags().StringVarP(&timeInterval, "time", "t", "-1h", "Example: -1h5m+5m")
	showCommand.Flags().StringVarP(&fields, "flags", "f", "", "Example: 1,2,3")

	return showCommand
}
