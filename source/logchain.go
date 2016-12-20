package source

import (
	"io"

	"github.com/kbence/logan/types"
)

// LogChain describes an interface to read logs from different files
type LogChain interface {
	Between(interval *types.TimeInterval) io.Reader
}
