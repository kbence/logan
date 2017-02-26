package pipeline

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/kbence/logan/types"
)

func TestSeek(t *testing.T) {
	var err error
	location, _ := time.LoadLocation("Local")
	startTime := time.Date(2017, 2, 26, 8, 0, 0, 0, location)
	endTime := time.Date(2017, 2, 26, 9, 0, 0, 0, location)
	reader := strings.NewReader("2017-02-26 08:00:05.123 This is a log line\n")
	expectedChunks := []string{"2017-02-26 08:00", ":05.123 This is ", "a log line\n"}

	buffer := make([]byte, 16)

	bufferedReader := NewTimeAwareBufferedReader(reader, types.NewTimeInterval(startTime, endTime))

	for idx, expected := range expectedChunks {
		var l int
		l, err = bufferedReader.Read(buffer)

		if idx < 2 && err != nil {
			t.Errorf("Error encountered for read #%d: '%s'", idx, err)
			return
		}

		if string(buffer[0:l]) != expected {
			t.Errorf("Expected to read '%s' for read #%d but reader returned '%s'!",
				expected, idx+1, string(buffer[0:l]))
			return
		}
	}

	if err != io.EOF {
		t.Errorf("Expected error io.EOF at end, got '%s' instead!", err)
	}
}
