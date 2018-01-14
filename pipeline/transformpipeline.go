package pipeline

import (
	"fmt"
	"strconv"

	"github.com/kbence/logan/types"
)

type TransformPipeline struct {
	inputChannel   types.LogLineChannel
	selectedFields []*types.IntInterval
}

func NewTransformPipeline(input types.LogLineChannel, fields []*types.IntInterval) *TransformPipeline {
	return &TransformPipeline{inputChannel: input, selectedFields: fields}
}

func createFieldList(intervals []*types.IntInterval, columns map[string]string) []int {
	fieldList := []int{}
	maxFieldID := 1

	for fieldIDStr := range columns {
		fieldID, err := strconv.ParseInt(fieldIDStr, 10, 32)

		if err == nil && maxFieldID < int(fieldID) {
			maxFieldID = int(fieldID)
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

		newLine := &types.LogLine{Line: line.Line, Date: line.Date, Columns: map[string]string{}}

		for idx, f := range createFieldList(fields, line.Columns) {
			newLine.Columns[fmt.Sprint(idx)] = line.Columns[fmt.Sprint(f)]
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
