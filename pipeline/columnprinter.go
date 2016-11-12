package pipeline

import (
	"fmt"
	"sort"

	"github.com/kbence/logan/types"
)

func getKeys(m types.ColumnList) []int {
	keys := []int{}

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func printColumnsInOrder(columns types.ColumnList) {
	format := "%s"
	keys := getKeys(columns)
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf(format, columns[key])
		format = " %s"
	}

	fmt.Printf("\n")
}
