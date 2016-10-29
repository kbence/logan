package filter

import (
	"testing"

	"github.com/kbence/logan/parser"
)

func TestSimpleEqualFilterWorks(t *testing.T) {
	filter := NewColumnFilter("$3 == \"teststring\"")
	testLine1 := &parser.LogLine{Columns: []string{"", "", "teststring"}}
	testLine2 := &parser.LogLine{Columns: []string{"", "", "string"}}

	if !filter.Match(testLine1) {
		t.Errorf("testLine1 (%s) should match filter (%s)", testLine1, filter)
	}

	if filter.Match(testLine2) {
		t.Errorf("testLine2 (%s) should not match filter (%s)", testLine1, filter)
	}
}

func TestShouldNotPanicIfColumnIndexIsWrong(t *testing.T) {
	filter := NewColumnFilter("$4 == \"teststring\"")
	testLine := &parser.LogLine{Columns: []string{"", "", "teststring"}}

	filter.Match(testLine)
}

func TestSimpleNotEqualsFilterWorks(t *testing.T) {
	filter := NewColumnFilter("$3 != \"teststring\"")
	testLine1 := &parser.LogLine{Columns: []string{"", "", "teststring"}}
	testLine2 := &parser.LogLine{Columns: []string{"", "", "string"}}

	if filter.Match(testLine1) {
		t.Errorf("testLine1 (%s) should match filter (%s)", testLine1, filter)
	}

	if !filter.Match(testLine2) {
		t.Errorf("testLine2 (%s) should not match filter (%s)", testLine1, filter)
	}
}
