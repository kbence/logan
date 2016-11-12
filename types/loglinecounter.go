package types

import "sort"

type UniqueLogLineCount struct {
	Line  *LogLine
	Count uint64
}

type uniqueLogLineSortFunc func(x, y *UniqueLogLineCount) bool

type uniqueLogLineCountSorter struct {
	counts      []*UniqueLogLineCount
	compareFunc uniqueLogLineSortFunc
}

func (s *uniqueLogLineCountSorter) Len() int {
	return len(s.counts)
}

func (s *uniqueLogLineCountSorter) Less(i, j int) bool {
	return s.compareFunc(s.counts[i], s.counts[j])
}

func (s *uniqueLogLineCountSorter) Swap(i, j int) {
	s.counts[i], s.counts[j] = s.counts[j], s.counts[i]
}

func sortByCountAsc(x, y *UniqueLogLineCount) bool {
	return x.Count < y.Count
}

func sortByCountDesc(x, y *UniqueLogLineCount) bool {
	return x.Count > y.Count
}

type LogLineCounter struct {
	counts map[string]*UniqueLogLineCount
}

func NewLogLineCounter() *LogLineCounter {
	return &LogLineCounter{counts: make(map[string]*UniqueLogLineCount)}
}

func (c *LogLineCounter) Add(line *LogLine) {
	hash := line.Columns.Join()
	counter, found := c.counts[hash]

	if !found {
		c.counts[hash] = &UniqueLogLineCount{Line: &LogLine{Columns: line.Columns}, Count: 1}
	} else {
		counter.Count++
	}
}

func (c *LogLineCounter) Max() uint64 {
	max := uint64(0)

	for _, line := range c.counts {
		if line.Count > max {
			max = line.Count
		}
	}

	return max
}

func (c *LogLineCounter) UniqueLines(sortDesc bool) []*UniqueLogLineCount {
	lines := []*UniqueLogLineCount{}

	for _, line := range c.counts {
		lines = append(lines, line)
	}

	var sortFunc uniqueLogLineSortFunc = sortByCountAsc

	if sortDesc {
		sortFunc = sortByCountDesc
	}

	sort.Sort(&uniqueLogLineCountSorter{counts: lines, compareFunc: sortFunc})

	return lines
}
