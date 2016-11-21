package main

import (
	"context"
	"flag"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
	"time"

	"github.com/google/subcommands"
	"github.com/kbence/logan/command"
	"github.com/kbence/logan/config"
	_ "github.com/kbence/logan/source"
)

var traceFilename = flag.String("trace", "", "Save Go trace data to file")
var traceFile *os.File

func startTrace() {
	if *traceFilename != "" {
		var err error
		traceFile, err = os.Create(*traceFilename)

		if err != nil {
			log.Panicf("ERROR opening trace file: %s", err)
		}

		trace.Start(traceFile)
	}
}

func stopTrace() {
	if traceFile != nil {
		trace.Stop()
		time.Sleep(time.Second)
		defer traceFile.Close()
	}
}

func main() {
	config := config.Load()

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(command.ListCommand(config), "")
	subcommands.Register(command.ShowCommand(config), "")
	subcommands.Register(command.UniqueCommand(config), "")
	subcommands.Register(command.InspectCommand(config), "")
	subcommands.Register(command.PlotCommand(config), "")

	flag.Parse()

	startTrace()
	exitCode := int(subcommands.Execute(context.Background()))
	stopTrace()

	os.Exit(exitCode)
}
