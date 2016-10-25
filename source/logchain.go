package source

import (
	"io"
	"time"
)

// LogChain describes an interface to read logs from different files
type LogChain interface {
	Last() io.Reader
	Between(start time.Time, end time.Time) io.Reader
}
