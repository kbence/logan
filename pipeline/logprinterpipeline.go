package pipeline

import "github.com/kbence/logan/types"

type LogPrinterPipeline struct {
	input types.LogLineChannel
}

func NewLogPrinterPipeline(input types.LogLineChannel) *LogPrinterPipeline {
	return &LogPrinterPipeline{input: input}
}

func (p *LogPrinterPipeline) Start() chan bool {
	exitChannel := make(chan bool)

	go func() {
		for {
			line, more := <-p.input

			if !more {
				break
			}

			printColumnsInOrder(line.Columns)
		}

		exitChannel <- true
	}()

	return exitChannel
}
