package source

import (
	"io/ioutil"
	"log"
	"sort"

	"github.com/kbence/logan/config"
)

func init() {
	logSourceFactories["scribe"] = func(cfg *config.Configuration) LogSource {
		return NewScribeLogSource(cfg.Scribe.Dirs)
	}
}

// ScribeLogSource implements source for Scribe logs
type ScribeLogSource struct {
	directories []string
}

// NewScribeLogSource returns a new instance of ScribeLogSource
func NewScribeLogSource(directories []string) *ScribeLogSource {
	return &ScribeLogSource{directories: directories}
}

// GetCategories returns directory names from scribe dirs
func (s *ScribeLogSource) GetCategories() []string {
	categoryMap := s.getCategoryDirectoryMap()
	categories := []string{}

	for cat := range categoryMap {
		categories = append(categories, cat)
	}

	sort.Strings(categories)

	return categories
}

func (s *ScribeLogSource) getCategoryDirectoryMap() map[string]string {
	categoryMap := make(map[string]string)

	for _, dir := range s.directories {
		files, err := ioutil.ReadDir(dir)

		if err != nil {
			log.Panicf("ERROR: %s\n", err)
		}

		for _, entry := range files {
			if entry.IsDir() {
				categoryMap[entry.Name()] = dir
			}
		}
	}

	return categoryMap
}

// GetChain returns a log chain for the given Scribe category
func (s *ScribeLogSource) GetChain(category string) LogChain {
	directory, found := s.getCategoryDirectoryMap()[category]

	if !found {
		return nil
	}

	return NewScribeLogChain(directory, category)
}
