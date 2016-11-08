package pipeline

import (
	"fmt"
	"sort"

	"github.com/kbence/logan/parser"
)

type OutputPipeline struct {
	input parser.LogLineChannel
}

func NewOutputPipeline(input parser.LogLineChannel) *OutputPipeline {
	return &OutputPipeline{input: input}
}

func getKeys(m map[int]string) []int {
	keys := []int{}

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func printColumnsInOrder(columns map[int]string) {
	format := "%s"
	keys := getKeys(columns)
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf(format, columns[key])
		format = " %s"
	}

	fmt.Printf("\n")
}

func (p *OutputPipeline) Start() {
	go func() {
		for {
			line, more := <-p.input

			if !more {
				break
			}

			printColumnsInOrder(line.Columns)
		}
	}()
}
