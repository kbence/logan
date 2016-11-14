package types

import (
	"bytes"
	"strings"
)

type ChartRenderer struct {
	Width  int
	Height int
	data   []uint64
}

func NewChartRenderer(width, height int) *ChartRenderer {
	return &ChartRenderer{Width: width, Height: height}
}

func (r *ChartRenderer) AddDataLine(data []uint64) {
	r.data = data
}

func max(data []uint64) uint64 {
	max := uint64(0)

	for _, value := range data {
		if value > max {
			max = value
		}
	}

	return max
}

func createBitmap(width, height int) [][]uint8 {
	m := make([][]uint8, height)

	for i := range m {
		m[i] = make([]uint8, width)
	}

	return m
}

func getRunes(s string) []rune {
	runes := []rune{}

	for _, rune := range s {
		runes = append(runes, rune)
	}

	return runes
}

func (r *ChartRenderer) Render() string {
	runes := getRunes(" ▘▝▀▖▌▞▛▗▚▐▜▄▙▟█")

	doubleWidth := r.Width * 2
	doubleHeight := r.Height * 2
	bitmap := createBitmap(doubleWidth, doubleHeight)
	max := max(r.data)

	if max == 0 {
		max = 1
	}

	for y, row := range bitmap {
		for x := range row {
			if x > 0 {
				y1 := uint64(doubleHeight) - ((r.data[x-1]+r.data[x])*uint64(doubleHeight)/2)/max
				y2 := uint64(doubleHeight) - r.data[x]*uint64(doubleHeight)/max

				if y2 < y1 {
					y1, y2 = y2, y1
				}

				if uint64(y) > y1 && uint64(y) < y2 || uint64(y) == y1 {
					bitmap[y][x] = 1
				}
			}

			if x < len(row)-1 {
				y1 := uint64(doubleHeight) - r.data[x]*uint64(doubleHeight)/max
				y2 := uint64(doubleHeight) - ((r.data[x]+r.data[x+1])*uint64(doubleHeight)/2)/max

				if y2 < y1 {
					y1, y2 = y2, y1
				}

				if uint64(y) >= y1 && uint64(y) < y2 || uint64(y) == y1 {
					bitmap[y][x] = 1
				}
			}
		}
	}

	lines := []string{}

	for y := 0; y < doubleHeight; y += 2 {
		line := bytes.NewBufferString("")

		for x := 0; x < doubleWidth; x += 2 {
			bits := bitmap[y][x] + bitmap[y][x+1]<<1 + bitmap[y+1][x]<<2 + bitmap[y+1][x+1]<<3
			line.WriteRune(runes[bits])
		}

		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}
