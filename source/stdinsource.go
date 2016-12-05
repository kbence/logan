package source

import "github.com/kbence/logan/config"

func init() {
	logSourceFactories["stdin"] = func(cfg *config.Configuration) LogSource {
		return NewStdinSource()
	}
}

type StdinSource struct{}

func NewStdinSource() *StdinSource {
	return &StdinSource{}
}

func (s *StdinSource) GetCategories() []string {
	return []string{"-"}
}

func (s *StdinSource) ContainsCategory(category string) bool {
	return category == "-"
}

func (s *StdinSource) GetChain(category string) LogChain {
	return NewStdinChain(category)
}
