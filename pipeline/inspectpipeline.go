package pipeline

import (
	"fmt"

	"github.com/kbence/logan/types"
)

type InspectPipeline struct {
	input types.LogLineChannel
}

func NewInspectPipeline(input types.LogLineChannel) *InspectPipeline {
	return &InspectPipeline{input: input}
}

func (p *InspectPipeline) printInspector(line *types.LogLine, id int) {
	if id > 0 {
		fmt.Println("----------------")
	}

	fmt.Printf("Line: %s\n", line.Line)

	for _, key := range line.Columns.SortedKeys() {
		fmt.Printf("%4s: %s\n", fmt.Sprintf("%d", key+1), line.Columns[key])
	}
}

func (p *InspectPipeline) Start() chan bool {
	exitChannel := make(chan bool)

	go func() {
		printedLines := 0
		for {
			line, more := <-p.input

			if !more {
				break
			}

			p.printInspector(line, printedLines)
			printedLines++

			if printedLines >= 5 {
				break
			}
		}

		exitChannel <- true
	}()

	return exitChannel
}
