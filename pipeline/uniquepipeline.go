package pipeline

import (
	"fmt"
	"strings"

	"github.com/kbence/logan/types"
	"github.com/kbence/logan/utils"
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
	timer            *utils.UpdateTimer
	lastPrintedLines int
}

func NewUniquePipeline(input types.LogLineChannel, settings UniqueSettings) *UniquePipeline {
	return &UniquePipeline{
		input:    input,
		counter:  types.NewLogLineCounter(),
		settings: settings,
		timer:    &utils.UpdateTimer{},
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

	for _, line := range p.counter.UniqueLines(true) {
		if p.settings.TopLimit > 0 && linesPrinted >= p.settings.TopLimit {
			return linesPrinted
		}

		fmt.Printf("%9d▕%s ", line.Count, drawPercentageBar(16, float64(line.Count)/float64(max)))
		printColumnsInOrder(line.Line.Columns)
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

func (p *UniquePipeline) updateIfNeeded() {
	if p.settings.TopLimit <= 0 {
		return
	}

	if p.timer.IsUpdateNeeded() {
		p.clearLines(p.lastPrintedLines, p.settings.TerminalWidth)

		numPrintedLines := p.printUniqueLines()

		if numPrintedLines > 0 {
			fmt.Print(terminfo.GoUpBy(numPrintedLines))
		}

		p.lastPrintedLines = numPrintedLines
	}
}

func (p *UniquePipeline) Start() chan bool {
	exitChannel := make(chan bool)

	go func() {
		for {
			line, more := <-p.input

			if !more {
				break
			}

			p.storeUniqueLine(line)
			p.updateIfNeeded()
		}

		p.printUniqueLines()
		exitChannel <- true
	}()

	return exitChannel
}
