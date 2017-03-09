package main

import (
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/kbence/logan/command"
	"github.com/kbence/logan/config"
	_ "github.com/kbence/logan/source"
)

func main() {
	config := config.Load()
	cmd := command.NewLoganCommand(config)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
