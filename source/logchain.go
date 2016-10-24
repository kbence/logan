package source

import "io"

// LogChain describes an interface to read logs from different files
type LogChain interface {
	Last() io.Reader
}
