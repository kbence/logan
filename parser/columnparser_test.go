package parser

import (
	"log"
	"testing"
	"time"
)

func makeChan(line string) LogLineChannel {
	ch := make(LogLineChannel)
	location, err := time.LoadLocation("Europe/Budapest")
	testDate := time.Date(2001, 2, 3, 4, 5, 6, 7, location)

	if err != nil {
		log.Panicf("ERROR: %s", location)
	}

	go func() {
		ch <- &LogLine{Date: testDate, Line: line}
		close(ch)
	}()

	return ch
}

func TestSimpleColumns(t *testing.T) {
	input := makeChan("first second third")
	output := make(LogLineChannel)
	defer close(output)

	go ParseColumns(input, output)

	result := <-output

	if len(result.Columns) != 3 {
		t.Error("Result must have 3 columns")
	} else if result.Columns[0] != "first" {
		t.Errorf("Column[0] \"%s\" != \"first\"", result.Columns[0])
	} else if result.Columns[1] != "second" {
		t.Errorf("Column[0] \"%s\" != \"first\"", result.Columns[1])
	} else if result.Columns[2] != "third" {
		t.Errorf("Column[0] \"%s\" != \"first\"", result.Columns[2])
	}
}

func TestQuotedColumns(t *testing.T) {
	input := makeChan("first \"second column\" \"third\"")
	output := make(LogLineChannel)
	defer close(output)

	go ParseColumns(input, output)

	result := <-output

	if len(result.Columns) != 3 {
		t.Error("Result must have 3 columns")
	} else if result.Columns[0] != "first" {
		t.Errorf("Column[0] \"%s\" != \"first\"", result.Columns[0])
	} else if result.Columns[1] != "second column" {
		t.Errorf("Column[0] \"%s\" != \"first\"", result.Columns[1])
	} else if result.Columns[2] != "third" {
		t.Errorf("Column[0] \"%s\" != \"first\"", result.Columns[2])
	}
}
