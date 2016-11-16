package pipeline

import "github.com/kbence/logan/types"

type TransformPipeline struct {
	inputChannel   types.LogLineChannel
	selectedFields []*types.IntInterval
}

func NewTransformPipeline(input types.LogLineChannel, fields []*types.IntInterval) *TransformPipeline {
	return &TransformPipeline{inputChannel: input, selectedFields: fields}
}

func createFieldList(intervals []*types.IntInterval, columns map[int]string) []int {
	fieldList := []int{}
	maxFieldID := 1

	for fieldID := range columns {
		if maxFieldID < fieldID {
			maxFieldID = fieldID
		}
	}

	for _, interval := range intervals {
		for fieldID := interval.Start; fieldID <= interval.End && fieldID <= maxFieldID; fieldID++ {
			fieldList = append(fieldList, fieldID)
		}
	}

	return fieldList
}

func selectFields(output types.LogLineChannel, input types.LogLineChannel, fields []*types.IntInterval) {
	for {
		line, more := <-input

		if !more {
			break
		}

		newLine := &types.LogLine{Line: line.Line, Date: line.Date, Columns: map[int]string{}}

		for idx, f := range createFieldList(fields, line.Columns) {
			newLine.Columns[idx] = line.Columns[f]
		}

		output <- newLine
	}

	close(output)
}

func (p *TransformPipeline) Start() types.LogLineChannel {
	outputChannel := types.NewLogLineChannel()

	go selectFields(outputChannel, p.inputChannel, p.selectedFields)

	return outputChannel
}
