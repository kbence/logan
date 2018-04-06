package natsort

import (
	"sort"
	"testing"

	"github.com/kbence/logan/utils"
)

func TestChunkify(t *testing.T) {
	cases := []struct {
		s        string
		expected chunks
	}{
		{"", chunks{}},
		{"test", chunks{newAlphaChunk("test")}},
		{"test1", chunks{newAlphaChunk("test"), newNumericChunk(1)}},
		{"three3chunks", chunks{newAlphaChunk("three"), newNumericChunk(3), newAlphaChunk("chunks")}},
		{"1test", chunks{newAlphaChunk(""), newNumericChunk(1), newAlphaChunk("test")}},
	}

	for _, c := range cases {
		chunks := chunkify(c.s)
		if !chunks.EqualsTo(c.expected) {
			t.Errorf("The result of chunkify(\"%s\") should be %s not %s",
				c.s, c.expected, chunks)
		}
	}
}

func TestNatCompare(t *testing.T) {
	cases := []struct {
		a, b     string
		expected int
	}{
		{"", "", 0},
		{"same", "same", 0},
		{"", "a", -1},
		{"test", "testabcd", -1},
		{"testabcd", "test", 1},
		{"test1", "testa", -1},
		{"test10", "test2", 1},
		{"test2", "test10", -1},
		{"test", "test3chunks", -1},
		{"test3chunks", "test", 1},
	}

	for _, c := range cases {
		result := NatCompare(c.a, c.b)
		if result != c.expected {
			t.Errorf("The result of comparation of '%s' and '%s' should be %d not %d",
				c.a, c.b, c.expected, result)
		}
	}
}

type SortTest struct {
	Items    []string
	Expected []string
}

func TestNatSort(t *testing.T) {
	cases := []SortTest{
		{
			[]string{"test item 1", "test item 10", "test item 2", "another test item"},
			[]string{"another test item", "test item 1", "test item 2", "test item 10"},
		},
		{
			[]string{"testslenghty", "tests", "test"},
			[]string{"test", "tests", "testslenghty"},
		},
	}

	for _, c := range cases {
		items := Strings(c.Items)
		expected := Strings(c.Expected)
		sort.Sort(items)

		if !utils.StringListsEqual(items, expected) {
			t.Errorf("Expected list: %s, got: %s", expected, items)
		}
	}
}

func TestNatSortReverse(t *testing.T) {
	cases := []SortTest{
		{
			[]string{"test item 1", "test item 10", "test item 2", "another test item"},
			[]string{"test item 10", "test item 2", "test item 1", "another test item"},
		},
		{
			[]string{"test", "tests", "testslenghty"},
			[]string{"testslenghty", "tests", "test"},
		},
	}

	for _, c := range cases {
		items := StringsReverse(c.Items)
		expected := StringsReverse(c.Expected)
		sort.Sort(items)

		if !utils.StringListsEqual(items, expected) {
			t.Errorf("Expected list: %s, got: %s", expected, items)
		}
	}
}
