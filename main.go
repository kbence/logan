package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

var configuration struct {
	Scribe struct {
		Dirs []string
	}
}

func loadConfig() {
	cfg, err := ini.LooseLoad(
		"/etc/logan.conf",
		fmt.Sprintf("%s/.logan.conf", os.Getenv("HOME")))

	if err != nil {
		log.Panicf("ERROR: %s\n", err)
	}

	configuration.Scribe.Dirs = []string{}
	scribeDirs := cfg.Section("scribe").Key("dirs").MustString("/mnt/scribe:/var/log/scribe")

	for _, dir := range strings.Split(scribeDirs, ":") {
		configuration.Scribe.Dirs = append(configuration.Scribe.Dirs, dir)
	}
}

func main() {
	loadConfig()

	logs := NewScribeLogSource(configuration.Scribe.Dirs)

	for _, cat := range logs.GetCategories() {
		fmt.Printf("Log category: %s\n", cat)
	}
}
