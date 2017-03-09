package command

import (
	"log"
	"time"

	"github.com/kbence/logan/config"
	"github.com/kbence/logan/pipeline"
	"github.com/kbence/logan/utils"
	"github.com/spf13/cobra"
)

// NewInspectCommand returns an instance of `inspect` command that is designed
// to inspect the first 5 lines from the selected log category that matches
// the current time and filter specification
func NewInspectCommand(cfg *config.Configuration) *cobra.Command {
	var timeInterval string

	inspectCommand := &cobra.Command{
		Use:   "inspect",
		Short: "Inspects the first 5 lines of log",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("You have to pass a log source to this command!")
			}

			p := pipeline.NewPipelineBuilder(pipeline.PipelineSettings{
				Category: args[0],
				Interval: utils.ParseTimeInterval(timeInterval, time.Now()),
				Filters:  args[1:],
				Fields:   utils.ParseIntervals(""),
				Config:   cfg,
				Output:   pipeline.OutputTypeInspector})
			p.Execute()
		},
	}

	inspectCommand.Flags().StringVarP(&timeInterval, "time", "t", "-1h", "Example: -1h5m+5m")

	return inspectCommand
}
