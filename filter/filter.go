package filter

import "github.com/kbence/logan/types"

type Filter interface {
	Match(line *types.LogLine) bool
}
