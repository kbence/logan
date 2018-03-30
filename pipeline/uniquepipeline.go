package pipeline

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/kbence/logan/types"
	"github.com/kbence/logan/utils/terminfo"
)

type UniqueSettings struct {
	TopLimit      int
	TerminalWidth int
}

type UniquePipeline struct {
	input            types.LogLineChannel
	counter          *types.LogLineCounter
	settings         UniqueSettings
	renderLock       *sync.Mutex
	lastPrintedLines int
}

func NewUniquePipeline(input types.LogLineChannel, settings UniqueSettings) *UniquePipeline {
	return &UniquePipeline{
		input:      input,
		counter:    types.NewLogLineCounter(),
		settings:   settings,
		renderLock: &sync.Mutex{},
	}
}

func (p *UniquePipeline) storeUniqueLine(line *types.LogLine) {
	p.counter.Add(line)
}

func drawPercentageBar(width int, percentage float64) string {
	barElements := []rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}
	barElementsLen := len(barElements)
	last := barElementsLen - 1
	bar := ""

	max := width * barElementsLen
	state := int(percentage * float64(max))

	for c := 0; c < width; c++ {
		var ch rune

		if (c+1)*barElementsLen < state {
			ch = barElements[last]
		} else if c*len(barElements) >= state {
			ch = barElements[0]
		} else {
			ch = barElements[(state-c*barElementsLen)-1]
		}

		bar = fmt.Sprintf("%s%c", bar, ch)
	}

	return bar
}

func (p *UniquePipeline) printUniqueLines() int {
	max := p.counter.Max()
	linesPrinted := 0
	truncateLines := p.settings.TopLimit > 0

	for _, line := range p.counter.UniqueLines(true) {
		if p.settings.TopLimit > 0 && linesPrinted >= p.settings.TopLimit {
			return linesPrinted
		}

		prefix := fmt.Sprintf("%9d▕%s ", line.Count, drawPercentageBar(16, float64(line.Count)/float64(max)))
		prefixLen := utf8.RuneCount([]byte(prefix))
		fmt.Print(prefix)

		if truncateLines {
			printColumnsInOrderWithLimit(line.Line.Columns, p.settings.TerminalWidth-prefixLen-1)
		} else {
			printColumnsInOrder(line.Line.Columns)
		}

		linesPrinted++
	}

	return linesPrinted
}

func (p *UniquePipeline) clearLines(lines, width int) {
	if lines <= 0 {
		return
	}

	deleteLine := strings.Repeat(" ", width)

	for i := 0; i < lines; i++ {
		fmt.Print(deleteLine)
	}

	fmt.Println(terminfo.GoUpBy(lines))
}

func (p *UniquePipeline) updateLoop(exitChannel chan bool) {
	for {
		select {
		case <-time.After(time.Second):
			p.renderLock.Lock()
			p.clearLines(p.lastPrintedLines, p.settings.TerminalWidth)

			numPrintedLines := p.printUniqueLines()

			if numPrintedLines > 0 {
				fmt.Print(terminfo.GoUpBy(numPrintedLines))
			}

			p.lastPrintedLines = numPrintedLines
			p.renderLock.Unlock()
			break

		case <-exitChannel:
			return
		}
	}
}

func (p *UniquePipeline) Start() chan bool {
	exitChannel := make(chan bool)
	updateExitChannel := make(chan bool)

	if p.settings.TopLimit > 0 {
		go p.updateLoop(updateExitChannel)
	}

	go func() {
		for {
			line, more := <-p.input

			if !more {
				break
			}

			p.renderLock.Lock()
			p.storeUniqueLine(line)
			p.renderLock.Unlock()
		}

		if p.settings.TopLimit > 0 {
			updateExitChannel <- true
		}

		p.renderLock.Lock()
		p.printUniqueLines()
		p.renderLock.Unlock()

		exitChannel <- true
	}()

	return exitChannel
}
