package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

// Configuration describes Logan's settings
type Configuration struct {
	Scribe struct {
		Dirs []string
	}
}

// Load tries to load configuration from several locations
func Load() *Configuration {
	var config Configuration

	cfg, err := ini.LooseLoad(
		"/etc/logan.conf",
		fmt.Sprintf("%s/.logan.conf", os.Getenv("HOME")))

	if err != nil {
		log.Panicf("ERROR: %s\n", err)
	}

	config.Scribe.Dirs = []string{}
	scribeDirs := cfg.Section("scribe").Key("dirs").MustString("/mnt/scribe:/var/log/scribe")

	for _, dir := range strings.Split(scribeDirs, ":") {
		config.Scribe.Dirs = append(config.Scribe.Dirs, dir)
	}

	return &config
}
