package natsort

import (
	"fmt"
	"strconv"

	"github.com/kbence/logan/utils"
)

type chunkType int

const (
	ctNumeric chunkType = iota
	ctAlpha
)

type chunk struct {
	Type   chunkType
	Number int
	Alpha  string
}

func (c chunk) String() string {
	if c.Type == ctAlpha {
		return fmt.Sprintf("\"%s\"", c.Alpha)
	}
	if c.Type == ctNumeric {
		return fmt.Sprintf("%d", c.Number)
	}

	return "nil"
}

type chunks []chunk

func (c chunks) EqualsTo(other chunks) bool {
	if c == nil && other == nil {
		return true
	}

	if c == nil || other == nil {
		return false
	}

	if len(c) != len(other) {
		return false
	}

	for idx, value := range c {
		if value != other[idx] {
			return false
		}
	}

	return true
}

func (c chunks) String() string {
	strings := []string{}

	for _, chunk := range c {
		strings = append(strings, chunk.String())
	}

	return utils.FormatStringArray(strings)
}

func newNumericChunk(number int) chunk {
	return chunk{Type: ctNumeric, Number: number}
}

func newAlphaChunk(alpha string) chunk {
	return chunk{Type: ctAlpha, Alpha: alpha}
}

func newChunk(cType chunkType, data []rune) chunk {
	var chk chunk

	switch cType {
	case ctAlpha:
		chk = newAlphaChunk(string(data))

	case ctNumeric:
		number, _ := strconv.Atoi(string(data))
		chk = newNumericChunk(number)
	}

	return chk
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func chunkify(s string) chunks {
	if len(s) > 0 {
		chk := chunks{}
		chunkType := ctAlpha
		buffer := []rune{}

		for _, r := range []rune(s) {
			if isDigit(r) {
				if chunkType == ctAlpha {
					chk = append(chk, newChunk(ctAlpha, buffer))
					buffer = []rune{}
				}

				chunkType = ctNumeric
			} else {
				if chunkType == ctNumeric {
					chk = append(chk, newChunk(ctNumeric, buffer))
					buffer = []rune{}
				}

				chunkType = ctAlpha
			}

			buffer = append(buffer, r)
		}

		if len(buffer) > 0 {
			if isDigit(buffer[0]) {
				chk = append(chk, newChunk(ctNumeric, buffer))
			} else {
				chk = append(chk, newChunk(ctAlpha, buffer))
			}
		}

		return chk
	}

	return chunks{}
}

func stringCompare(a, b string) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}

func intCompare(a, b int) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}

func NatCompare(a, b string) int {
	aChunks, bChunks := chunkify(a), chunkify(b)
	aLen, bLen := len(aChunks), len(bChunks)
	minLen := utils.MinInt(aLen, bLen)

	for i := 0; i < minLen; i++ {
		if aChunks[i].Type == ctAlpha {
			if cmp := stringCompare(aChunks[i].Alpha, bChunks[i].Alpha); cmp != 0 {
				return cmp
			}
		} else {
			if cmp := intCompare(aChunks[i].Number, bChunks[i].Number); cmp != 0 {
				return cmp
			}
		}
	}

	return intCompare(aLen, bLen)
}

type Strings []string
type StringsReverse []string

func (s Strings) Len() int           { return len(s) }
func (s Strings) Swap(a, b int)      { s[a], s[b] = s[b], s[a] }
func (s Strings) Less(a, b int) bool { return NatCompare(s[a], s[b]) < 0 }
func (s Strings) String() string     { return utils.FormatStringArray(s) }

func (s StringsReverse) Len() int           { return len(s) }
func (s StringsReverse) Swap(a, b int)      { s[a], s[b] = s[b], s[a] }
func (s StringsReverse) Less(a, b int) bool { return NatCompare(s[a], s[b]) > 0 }
func (s StringsReverse) String() string     { return utils.FormatStringArray(s) }
