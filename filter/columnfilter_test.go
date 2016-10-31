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
		t.Errorf("testLine2 (%s) should not match filter (%s)", testLine2, filter)
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
		t.Errorf("testLine2 (%s) should not match filter (%s)", testLine2, filter)
	}
}

func TestAndWorksBetweenFilterExpressions(t *testing.T) {
	filter := NewColumnFilter("$1 == \"1\" AND $3 == \"teststring\"")
	testLineBothMatches := &parser.LogLine{Columns: []string{"1", "", "teststring"}}
	testLineFirstMatches := &parser.LogLine{Columns: []string{"1", "", "string"}}
	testLineSecondMatches := &parser.LogLine{Columns: []string{"2", "", "teststring"}}
	testLineNothingMatches := &parser.LogLine{Columns: []string{"2", "", "string"}}

	if !filter.Match(testLineBothMatches) {
		t.Errorf("testLineBothMatches (%s) should match filter (%s)", testLineBothMatches, filter)
	}

	if filter.Match(testLineFirstMatches) {
		t.Errorf("testLineFirstMatches (%s) should not match filter (%s)", testLineFirstMatches, filter)
	}

	if filter.Match(testLineSecondMatches) {
		t.Errorf("testLineSecondMatches (%s) should not match filter (%s)", testLineSecondMatches, filter)
	}

	if filter.Match(testLineNothingMatches) {
		t.Errorf("testLineNothingMatches (%s) should not match filter (%s)", testLineNothingMatches, filter)
	}
}

func TestOrWorksBetweenFilterExpressions(t *testing.T) {
	filter := NewColumnFilter("$1 == \"1\" OR $3 == \"teststring\"")
	testLineBothMatches := &parser.LogLine{Columns: []string{"1", "", "teststring"}}
	testLineFirstMatches := &parser.LogLine{Columns: []string{"1", "", "string"}}
	testLineSecondMatches := &parser.LogLine{Columns: []string{"2", "", "teststring"}}
	testLineNothingMatches := &parser.LogLine{Columns: []string{"2", "", "string"}}

	if !filter.Match(testLineBothMatches) {
		t.Errorf("testLineBothMatches (%s) should match filter (%s)", testLineBothMatches, filter)
	}

	if !filter.Match(testLineFirstMatches) {
		t.Errorf("testLineFirstMatches (%s) should not match filter (%s)", testLineFirstMatches, filter)
	}

	if !filter.Match(testLineSecondMatches) {
		t.Errorf("testLineSecondMatches (%s) should not match filter (%s)", testLineSecondMatches, filter)
	}

	if filter.Match(testLineNothingMatches) {
		t.Errorf("testLineNothingMatches (%s) should not match filter (%s)", testLineNothingMatches, filter)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	filter := NewColumnFilter("$1 == \"1\" OR $2 == \"second\" AND $3 == \"teststring\"")
	testLineMatches1 := &parser.LogLine{Columns: []string{"1", "no match", "no match"}}
	testLineMatches2 := &parser.LogLine{Columns: []string{"no match", "second", "teststring"}}
	testLineMatches3 := &parser.LogLine{Columns: []string{"1", "second", "teststring"}}
	testLineNoMatch1 := &parser.LogLine{Columns: []string{"no match", "second", "no match"}}
	testLineNoMatch2 := &parser.LogLine{Columns: []string{"no match", "no match", "second"}}

	if !filter.Match(testLineMatches1) {
		t.Errorf("testLineMatches1 (%s) should match filter (%s)", testLineMatches1, filter)
	}

	if !filter.Match(testLineMatches2) {
		t.Errorf("testLineMatches2 (%s) should match filter (%s)", testLineMatches2, filter)
	}

	if !filter.Match(testLineMatches3) {
		t.Errorf("testLineMatches3 (%s) should match filter (%s)", testLineMatches3, filter)
	}

	if filter.Match(testLineNoMatch1) {
		t.Errorf("testLineNoMatch1 (%s) should not match filter (%s)", testLineNoMatch1, filter)
	}

	if filter.Match(testLineNoMatch2) {
		t.Errorf("testLineNoMatch2 (%s) should not match filter (%s)", testLineNoMatch2, filter)
	}
}
