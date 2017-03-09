package command

import (
	"log"
	"os"
	"runtime/trace"
	"time"

	"github.com/kbence/logan/config"
	"github.com/spf13/cobra"
)

// NewLoganCommand returns the root command that contains all
// logan commands and starts/stops go trace
func NewLoganCommand(cfg *config.Configuration) *cobra.Command {
	var traceFilename string
	var traceFile *os.File

	command := &cobra.Command{
		Use:   "logan",
		Short: "Command line tool for analyzing logs",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if traceFilename != "" {
				var err error
				traceFile, err = os.Create(traceFilename)

				if err != nil {
					log.Panicf("ERROR opening trace file: %s", err)
				}

				trace.Start(traceFile)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if traceFile != nil {
				trace.Stop()
				time.Sleep(time.Second)
				defer traceFile.Close()
			}
		},
	}

	command.PersistentFlags().StringVarP(&traceFilename, "trace", "", "", "Saves go trace to the specified file")
	command.AddCommand(NewListCommand(cfg))
	command.AddCommand(NewInspectCommand(cfg))
	command.AddCommand(NewShowCommand(cfg))
	command.AddCommand(NewPlotCommand(cfg))
	command.AddCommand(NewUniqCommand(cfg))

	return command
}
