package source

import (
	"io"

	"github.com/kbence/logan/types"
)

// LogChain describes an interface to read logs from different files
type LogChain interface {
	Last() io.Reader
	Between(interval *types.TimeInterval) io.Reader
}
