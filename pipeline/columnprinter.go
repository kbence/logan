package pipeline

import (
	"fmt"
	"sort"

	"github.com/kbence/logan/types"
)

func getKeys(m types.ColumnList) []string {
	keys := []string{}

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func getLinesFromColumn(columns types.ColumnList) string {
	line := ""
	format := "%s%s"
	keys := getKeys(columns)
	sort.Strings(keys)

	for _, key := range keys {
		line = fmt.Sprintf(format, line, columns[key])
		format = "%s %s"
	}

	return line
}

func printColumnsInOrder(columns types.ColumnList) {
	fmt.Printf("%s\n", getLinesFromColumn(columns))
}

func printColumnsInOrderWithLimit(columns types.ColumnList, limit int) {
	line := getLinesFromColumn(columns)

	if len(line) > limit {
		line = fmt.Sprintf("%s%s", line[0:limit-1], "…")
	}

	fmt.Printf("%s\n", line)
}
