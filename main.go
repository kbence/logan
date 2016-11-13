package main

import (
	"flag"
	"os"

	"golang.org/x/net/context"

	"github.com/google/subcommands"

	"github.com/kbence/logan/command"
	"github.com/kbence/logan/config"
	_ "github.com/kbence/logan/source"
)

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
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
