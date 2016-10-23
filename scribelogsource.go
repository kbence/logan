package main

import (
	"io/ioutil"
	"log"
	"sort"
)

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
	categoryMap := make(map[string]bool)

	for _, dir := range s.directories {
		files, err := ioutil.ReadDir(dir)

		if err != nil {
			log.Panicf("ERROR: %s\n", err)
		}

		for _, entry := range files {
			if entry.IsDir() {
				categoryMap[entry.Name()] = true
			}
		}
	}

	categories := []string{}

	for cat := range categoryMap {
		categories = append(categories, cat)
	}

	sort.Strings(categories)

	return categories
}
