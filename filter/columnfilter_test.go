package filter

import (
	"testing"

	"github.com/kbence/logan/types"
)

type expectation struct {
	Matches bool
	Line    *types.LogLine
}

func newExpectation(matches bool, columns ...string) *expectation {
	indexedColumns := types.ColumnList{}

	for i, col := range columns {
		indexedColumns[i+1] = col
	}

	return &expectation{Matches: matches, Line: &types.LogLine{Columns: indexedColumns}}
}

func expectMatch(columns ...string) *expectation {
	return newExpectation(true, columns...)
}

func expectDoesntMatch(columns ...string) *expectation {
	return newExpectation(false, columns...)
}

func testFilter(t *testing.T, filterString string, expectations ...*expectation) {
	filter := NewColumnFilter(filterString)

	for _, exp := range expectations {
		if exp.Matches {
			if !filter.Match(exp.Line) {
				t.Errorf("%s should match line %s", filter, exp.Line)
			}
		} else {
			if filter.Match(exp.Line) {
				t.Errorf("%s should not match line %s", filter, exp.Line)
			}
		}
	}
}

func TestSimpleEqualFilterWorks(t *testing.T) {
	testFilter(t, "$3 == \"teststring\"",
		expectMatch("", "", "teststring"),
		expectDoesntMatch("", "", "string"))

}

func TestSimpleEqualFilterWorksWithoutSpaces(t *testing.T) {
	testFilter(t, "$3==\"teststring\"",
		expectMatch("", "", "teststring"),
		expectDoesntMatch("", "", "string"))
}

func TestShouldNotPanicIfColumnIndexIsWrong(t *testing.T) {
	testFilter(t, "$4 == \"teststring\"",
		expectDoesntMatch("", "", "teststring"))
}

func TestSimpleNotEqualsFilterWorks(t *testing.T) {
	testFilter(t, "$3 != \"teststring\"",
		expectDoesntMatch("", "", "teststring"),
		expectMatch("", "", "string"))
}

func TestAndWorksBetweenFilterExpressions(t *testing.T) {
	testFilter(t, "$1 == \"1\" AND $3 == \"teststring\"",
		expectMatch("1", "", "teststring"),
		expectDoesntMatch("1", "", "string"),
		expectDoesntMatch("2", "", "teststring"),
		expectDoesntMatch("2", "", "string"))
}

func TestOrWorksBetweenFilterExpressions(t *testing.T) {
	testFilter(t, "$1 == \"1\" OR $3 == \"teststring\"",
		expectMatch("1", "", "teststring"),
		expectMatch("1", "", "string"),
		expectMatch("2", "", "teststring"),
		expectDoesntMatch("2", "", "string"))
}

func TestOperatorPrecedence(t *testing.T) {
	testFilter(t, "$1 == \"1\" OR $2 == \"second\" AND $3 == \"teststring\"",
		expectMatch("1", "no match", "no match"),
		expectMatch("no match", "second", "teststring"),
		expectMatch("1", "second", "teststring"),
		expectDoesntMatch("no match", "second", "no match"),
		expectDoesntMatch("no match", "no match", "second"))
}

func TestRegularExpressionWorks(t *testing.T) {
	testFilter(t, "$1 ~= \"^test\"",
		expectMatch("teststring"),
		expectDoesntMatch("another teststring"))

	testFilter(t, "\"interesting idea of reverse matching\" ~= $1",
		expectMatch("reverse matching$"),
		expectDoesntMatch("weird idea"))
}
