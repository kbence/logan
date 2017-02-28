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
	Generic struct {
		Dirs     []string
		MaxDepth int
	}
}

const (
	scribeSection  = "scribe"
	genericSection = "generic"
)

type iniFile ini.File

func (f *iniFile) extractDirs(section, defaultValue string) []string {
	dirSpec := (*ini.File)(f).Section(section).Key("dirs").MustString(defaultValue)
	dirs := []string{}

	for _, dir := range strings.Split(dirSpec, ":") {
		dirs = append(dirs, dir)
	}

	return dirs
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

	config.Scribe.Dirs = (*iniFile)(cfg).extractDirs(scribeSection, "/mnt/scribe:/var/log/scribe")
	config.Generic.Dirs = (*iniFile)(cfg).extractDirs(genericSection, "/var/log")
	config.Generic.MaxDepth = cfg.Section(genericSection).Key("recursion").MustInt(1)

	return &config
}
