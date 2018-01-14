package parser

import (
	"fmt"
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
	} else if result.Columns["1"] != "first" {
		t.Errorf("Column[1] \"%s\" != \"first\" '%s'", result.Columns["1"], result.Columns["1"])
	} else if result.Columns["2"] != "second" {
		t.Errorf("Column[2] \"%s\" != \"second\"", result.Columns["2"])
	} else if result.Columns["3"] != "third" {
		fmt.Println(result.Columns)
		t.Errorf("Column[3] \"%s\" != \"third\"", result.Columns["3"])
	}
}

func TestQuotedColumns(t *testing.T) {
	input := makeChan("first \"second column\" 'third'")
	output := make(types.LogLineChannel)

	go ParseColumns(output, input)

	result := <-output

	if len(result.Columns) != 3 {
		t.Error("Result must have 3 columns")
	} else if result.Columns["2"] != "second column" {
		t.Errorf("Column[2] \"%s\" != \"second column\"", result.Columns["2"])
	} else if result.Columns["3"] != "third" {
		t.Errorf("Column[3] \"%s\" != \"third\"", result.Columns["3"])
	}
}

func TestBracketedColumns(t *testing.T) {
	input := makeChan("first [second column] (third column)")
	output := make(types.LogLineChannel)

	go ParseColumns(output, input)

	result := <-output

	if len(result.Columns) != 3 {
		t.Error("Result must have 3 columns")
	} else if result.Columns["2"] != "[second column]" {
		t.Errorf("Column[2] \"%s\" != \"[second column]\"", result.Columns["2"])
	} else if result.Columns["3"] != "(third column)" {
		t.Errorf("Column[3] \"%s\" != \"third column\"", result.Columns["3"])
	}
}

func TestStructuredColumns(t *testing.T) {
	input := makeChan("{\"time\": \"1488888888\", \"message\": \"this is a structured log line\"}")
	output := make(types.LogLineChannel)

	go ParseColumns(output, input)

	result := <-output

	if len(result.Columns) != 2 {
		t.Errorf("Result must have 2 columns (have %d)!", len(result.Columns))
		fmt.Println(result.Columns)
	}
}

func TestIsJSONLine(t *testing.T) {
	cases := []struct {
		line     string
		expected bool
	}{
		{"", false},
		{"{", false},
		{"{}", true},
		{"{\"test\":1}", true},
		{"{\"test\": 2} ", true},
		{" {\"test\":3}", true},
		{"x{\"test\": 4}", false},
	}

	for _, c := range cases {
		if isJSONLine(c.line) != c.expected {
			not := ""
			if !c.expected {
				not = " not"
			}
			t.Errorf("Line '%s' should%s be parsed as JSON!", c.line, not)
		}
	}
}
