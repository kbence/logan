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
