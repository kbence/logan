package pipeline

import (
	"io"
	"time"

	"github.com/kbence/logan/parser"
	"github.com/kbence/logan/types"
)

const readBlockSize = 1024 * 1024

type TimeAwareBufferedReader struct {
	reader       io.Reader
	interval     *types.TimeInterval
	buffer       []byte
	bufferPos    int
	bufferLen    int
	startReached bool
	ended        bool
	endError     error
}

func NewTimeAwareBufferedReader(reader io.Reader, interval *types.TimeInterval) *TimeAwareBufferedReader {
	return &TimeAwareBufferedReader{reader: reader, interval: interval}
}

func (r *TimeAwareBufferedReader) readNext() error {
	var err error
	var length int

	if r.buffer == nil {
		r.buffer = make([]byte, readBlockSize)
	}

	r.bufferLen = 0

	for err == nil && r.bufferLen < readBlockSize {
		length, err = r.reader.Read(r.buffer[r.bufferLen:len(r.buffer)])
		r.bufferLen += length
	}

	if err != nil && r.bufferLen > 0 {
		r.endError = err
		err = nil
	}

	r.bufferPos = 0

	return err
}

func (r *TimeAwareBufferedReader) seek() error {
	for {
		err := r.readNext()

		if err != nil && err != io.EOF {
			r.ended = true
			r.endError = err
			break
		}

		var date *time.Time

		date = parseLastDate(r.buffer, r.bufferLen)

		if date != nil && !date.Before(r.interval.StartTime) {
			r.startReached = true
			return nil
		}
	}

	return nil
}

func (r *TimeAwareBufferedReader) Read(slice []byte) (int, error) {
	if !r.startReached {
		err := r.seek()

		if err != nil {
			return 0, err
		}
	}

	var err error
	lengthToRead := len(slice)
	currentSlicePos := 0

	for currentSlicePos < lengthToRead && err == nil {
		if r.bufferPos >= r.bufferLen {
			err = r.readNext()
		}

		remainingLength := r.bufferLen - r.bufferPos
		readFromBuffer := min(remainingLength, len(slice)-currentSlicePos)

		if readFromBuffer > 0 {
			copy(slice[currentSlicePos:currentSlicePos+readFromBuffer],
				r.buffer[r.bufferPos:r.bufferPos+readFromBuffer])

			currentSlicePos += readFromBuffer
			r.bufferPos += readFromBuffer
		}

		if date := parseLastDate(r.buffer, r.bufferPos); date != nil && date.After(r.interval.EndTime) {
			err = io.EOF
		} else if r.ended && r.bufferPos >= r.bufferLen {
			err = r.endError
		}
	}

	return currentSlicePos, err
}

func extractLastFullLine(buffer []byte, length int) []byte {
	lineBounds := []int{0, length}
	currentBound := 1

	for pos := length - 1; pos >= 0 && currentBound >= 0; pos-- {
		if buffer[pos] == '\n' {
			lineBounds[currentBound] = pos + 1
			currentBound--
		}
	}

	return buffer[lineBounds[0]:lineBounds[1]]
}

func parseLastDate(buffer []byte, length int) *time.Time {
	var date *time.Time
	end := length

	for date == nil && end > 0 {
		line := string(extractLastFullLine(buffer, length))
		date = parser.ParseDate(line)
		end -= len(line)
	}

	return date
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
