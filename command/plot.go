package command

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kbence/logan/config"
	"github.com/kbence/logan/pipeline"
	"github.com/kbence/logan/types"
	"github.com/kbence/logan/utils"
	"github.com/spf13/cobra"
)

// NewPlotCommand returns the command that starts up the pipeline for
// plotting a chart
func NewPlotCommand(cfg *config.Configuration) *cobra.Command {
	var timeInterval string
	var fields string
	var mode string
	var autoUpdate bool

	plotCommand := &cobra.Command{
		Use:   "plot",
		Short: "Plots a time-based line chart from the number of log lines",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("You have to pass a log source to this command!")
			}

			interval := utils.ParseTimeInterval(timeInterval, time.Now())
			width, height := utils.GetTerminalDimensions()

			p := pipeline.NewPipelineBuilder(pipeline.PipelineSettings{
				Category: args[0],
				Interval: interval,
				Filters:  args[1:],
				Fields:   utils.ParseIntervals(fields),
				Config:   cfg,
				Output:   pipeline.OutputTypeLineChart,
				OutputSettings: pipeline.LineChartSettings{
					Mode:            mode,
					Width:           width,
					Height:          height - 1,
					Interval:        interval,
					FrequentUpdates: autoUpdate}})
			p.Execute()
		},
	}

	plotCommand.Flags().StringVarP(&timeInterval, "time", "t", "-1h", "Example: -1h5m+5m")
	plotCommand.Flags().StringVarP(&fields, "fields", "f", "", "Example: 1,2,3")
	plotCommand.Flags().StringVarP(&mode, "mode", "m", "braille",
		fmt.Sprintf("One of the following modes: %s.", strings.Join(types.CharacterSets.GetNames(), ", ")))
	plotCommand.Flags().BoolVarP(&autoUpdate, "auto-update", "u", true, "Auto-update chart during log parsing")

	return plotCommand
}
