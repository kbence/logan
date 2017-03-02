package pipeline

import (
	"bufio"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/kbence/logan/types"
)

func TestSeekRunsToEnd(t *testing.T) {
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

func TestSeekFirstReadHitsEnd(t *testing.T) {
	var err error
	var length int
	location, _ := time.LoadLocation("Local")
	startTime := time.Date(2017, 2, 26, 8, 0, 0, 0, location)
	endTime := time.Date(2017, 2, 26, 9, 0, 0, 0, location)
	logContent := "2017-02-26 08:00:05.123 This is a log line"
	reader := strings.NewReader(logContent)

	buffer := make([]byte, 1024)

	bufferedReader := NewTimeAwareBufferedReader(reader, types.NewTimeInterval(startTime, endTime))
	length, err = bufferedReader.Read(buffer)

	if length != 42 {
		t.Errorf("Length is not matching (expected: 42, actual: %d)!", length)
		return
	}

	if string(buffer[0:length]) != logContent {
		t.Errorf("Expected to read '%s' but reader returned '%s'!",
			logContent, string(buffer[0:length]))
		return
	}

	if err != io.EOF {
		t.Errorf("Expected error io.EOF at end, got '%s' instead!", err)
	}
}

func TestSeekWithReadLine(t *testing.T) {
	var err error

	location, _ := time.LoadLocation("Local")
	startTime := time.Date(2017, 2, 26, 8, 0, 0, 0, location)
	endTime := time.Date(2017, 2, 26, 9, 0, 0, 0, location)
	logContent := "2017-02-26 08:00:05.123 This is a log line\n" +
		"2017-02-26 08:00:05.123 This is another log line"
	logReader := strings.NewReader(logContent)

	multiReader := io.MultiReader(logReader)
	bufferedReader := NewTimeAwareBufferedReader(multiReader, types.NewTimeInterval(startTime, endTime))
	reader := bufio.NewReader(bufferedReader)

	line1, err := reader.ReadString('\n')
	line2, err := reader.ReadString('\n')

	if line1 != logContent[0:43] {
		t.Errorf("Expected to read '%s' but reader returned '%s'!",
			logContent[0:43], string(line1))
		return
	}

	if err != nil {
		t.Errorf("Expected no error, got '%s' instead!", err)
	}

	if line2 != logContent[43:91] {
		t.Errorf("Expected to read '%s' but reader returned '%s'!",
			logContent[43:91], string(line2))
		return
	}

	if err != io.EOF {
		t.Errorf("Expected error io.EOF at end, got '%s' instead!", err)
	}
}
