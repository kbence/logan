package source

import "github.com/kbence/logan/config"

type logSourceFactory func(*config.Configuration) LogSource

var logSourceFactories = map[string]logSourceFactory{}

// GetLogSources returns all log sources in a map
func GetLogSources(cfg *config.Configuration) map[string]LogSource {
	sources := map[string]LogSource{}

	for name, factory := range logSourceFactories {
		sources[name] = factory(cfg)
	}

	return sources
}

// GetLogSource returns the source with the given name
func GetLogSource(cfg *config.Configuration, source string) LogSource {
	factory, found := logSourceFactories[source]

	if !found {
		return nil
	}

	sourceObject := factory(cfg)

	return sourceObject
}

func GetSourcesForCategory(cfg *config.Configuration, category string) []LogSource {
	sources := []LogSource{}

	for _, factory := range logSourceFactories {
		source := factory(cfg)
		if source.ContainsCategory(category) {
			sources = append(sources, source)
		}
	}

	return sources
}
