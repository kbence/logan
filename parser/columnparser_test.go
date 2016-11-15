package parser

import (
	"log"
	"testing"
	"time"

	"github.com/kbence/logan/types"
)

func makeChan(line string) types.LogLineChannel {
	ch := make(types.LogLineChannel)
	location, err := time.LoadLocation("Europe/Budapest")
	testDate := time.Date(2001, 2, 3, 4, 5, 6, 7, location)

	if err != nil {
		log.Panicf("ERROR: %s", location)
	}

	go func() {
		ch <- &types.LogLine{Date: testDate, Line: line}
		close(ch)
	}()

	return ch
}

func TestSimpleColumns(t *testing.T) {
	input := makeChan("first second third")
	output := make(types.LogLineChannel)

	go ParseColumns(output, input)

	result := <-output

	if len(result.Columns) != 3 {
		t.Error("Result must have 3 columns")
	} else if result.Columns[1] != "first" {
		t.Errorf("Column[1] \"%s\" != \"first\"", result.Columns[0])
	} else if result.Columns[2] != "second" {
		t.Errorf("Column[2] \"%s\" != \"second\"", result.Columns[1])
	} else if result.Columns[3] != "third" {
		t.Errorf("Column[3] \"%s\" != \"third\"", result.Columns[2])
	}
}

func TestQuotedColumns(t *testing.T) {
	input := makeChan("first \"second column\" \"third\"")
	output := make(types.LogLineChannel)

	go ParseColumns(output, input)

	result := <-output

	if len(result.Columns) != 3 {
		t.Error("Result must have 3 columns")
	} else if result.Columns[2] != "second column" {
		t.Errorf("Column[1] \"%s\" != \"second column\"", result.Columns[1])
	} else if result.Columns[3] != "third" {
		t.Errorf("Column[2] \"%s\" != \"third\"", result.Columns[2])
	}
}

func TestBracketedColumns(t *testing.T) {
	input := makeChan("first [second column] [third]")
	output := make(types.LogLineChannel)

	go ParseColumns(output, input)

	result := <-output

	if len(result.Columns) != 3 {
		t.Error("Result must have 3 columns")
	} else if result.Columns[2] != "second column" {
		t.Errorf("Column[1] \"%s\" != \"second column\"", result.Columns[1])
	} else if result.Columns[3] != "third" {
		t.Errorf("Column[2] \"%s\" != \"third\"", result.Columns[2])
	}
}
