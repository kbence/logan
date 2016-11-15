package types

import (
	"bytes"
	"strings"
)

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

type CharacterSet struct {
	name  string
	runes []rune
	bits  [][]uint8
}

func (s *CharacterSet) HorizontalMultiplier() int { return len(s.bits[0]) }
func (s *CharacterSet) VerticalMultiplier() int   { return len(s.bits) }

type CharacterSetList []CharacterSet

var CharacterSets = CharacterSetList{
	CharacterSet{
		name:  "classic",
		runes: getRunes(" '.|"),
		bits: [][]uint8{
			[]uint8{0},
			[]uint8{1},
		},
	},
	CharacterSet{
		name:  "block",
		runes: getRunes(" █"),
		bits: [][]uint8{
			[]uint8{0},
		},
	},
	CharacterSet{
		name:  "quad",
		runes: getRunes(" ▘▝▀▖▌▞▛▗▚▐▜▄▙▟█"),
		bits: [][]uint8{
			[]uint8{0, 1},
			[]uint8{2, 3},
		},
	},
	CharacterSet{
		name: "brailles",
		runes: getRunes("" +
			"⠀⠁⠂⠃⠄⠅⠆⠇⠈⠉⠊⠋⠌⠍⠎⠏" +
			"⠐⠑⠒⠓⠔⠕⠖⠗⠘⠙⠚⠛⠜⠝⠞⠟" +
			"⠠⠡⠢⠣⠤⠥⠦⠧⠨⠩⠪⠫⠬⠭⠮⠯" +
			"⠰⠱⠲⠳⠴⠵⠶⠷⠸⠹⠺⠻⠼⠽⠾⠿" +
			"⡀⡁⡂⡃⡄⡅⡆⡇⡈⡉⡊⡋⡌⡍⡎⡏" +
			"⡐⡑⡒⡓⡔⡕⡖⡗⡘⡙⡚⡛⡜⡝⡞⡟" +
			"⡠⡡⡢⡣⡤⡥⡦⡧⡨⡩⡪⡫⡬⡭⡮⡯" +
			"⡰⡱⡲⡳⡴⡵⡶⡷⡸⡹⡺⡻⡼⡽⡾⡿" +
			"⢀⢁⢂⢃⢄⢅⢆⢇⢈⢉⢊⢋⢌⢍⢎⢏" +
			"⢐⢑⢒⢓⢔⢕⢖⢗⢘⢙⢚⢛⢜⢝⢞⢟" +
			"⢠⢡⢢⢣⢤⢥⢦⢧⢨⢩⢪⢫⢬⢭⢮⢯" +
			"⢰⢱⢲⢳⢴⢵⢶⢷⢸⢹⢺⢻⢼⢽⢾⢿" +
			"⣀⣁⣂⣃⣄⣅⣆⣇⣈⣉⣊⣋⣌⣍⣎⣏" +
			"⣐⣑⣒⣓⣔⣕⣖⣗⣘⣙⣚⣛⣜⣝⣞⣟" +
			"⣠⣡⣢⣣⣤⣥⣦⣧⣨⣩⣪⣫⣬⣭⣮⣯" +
			"⣰⣱⣲⣳⣴⣵⣶⣷⣸⣹⣺⣻⣼⣽⣾⣿"),
		bits: [][]uint8{
			[]uint8{0, 3},
			[]uint8{1, 4},
			[]uint8{2, 5},
			[]uint8{6, 7},
		},
	},
}

func (s CharacterSetList) Select(name string) *CharacterSet {
	for _, charSet := range CharacterSets {
		if name == charSet.name {
			return &charSet
		}
	}

	return &CharacterSets[0]
}

func (s CharacterSetList) GetNames() []string {
	names := []string{}

	for _, charSet := range CharacterSets {
		names = append(names, charSet.name)
	}

	return names
}

type ChartSettings struct {
	Mode        string
	Width       int
	Height      int
	Border      bool
	XAxisLabels bool
	YAxisLabels bool
}

func (s *ChartSettings) EffectiveWidth() int {
	width := s.Width

	if s.Border {
		width--
	}

	return width

}

func (s *ChartSettings) EffectiveHeight() int {
	height := s.Height

	if s.Border {
		height--
	}

	return height
}

func (s *ChartSettings) SamplerSize() int {
	return s.EffectiveWidth() * CharacterSets.Select(s.Mode).HorizontalMultiplier()
}

type ChartRenderer struct {
	settings *ChartSettings
	data     []uint64
}

func NewChartRenderer(settings *ChartSettings) *ChartRenderer {
	return &ChartRenderer{settings: settings}
}

func (r *ChartRenderer) AddDataLine(data []uint64) {
	r.data = data
}

func (r *ChartRenderer) Render() string {
	set := CharacterSets.Select(r.settings.Mode)

	chartAreaWidth := r.settings.Width
	chartAreaHeight := r.settings.Height

	wMult := set.HorizontalMultiplier()
	hMult := set.VerticalMultiplier()

	mulWidth := uint64(chartAreaWidth * wMult)
	mulHeight := uint64(chartAreaHeight * hMult)
	bitmap := createBitmap(int(mulWidth), int(mulHeight))
	max := max(r.data)

	if max == 0 {
		max = 1
	}

	for y, row := range bitmap {
		for x := range row {
			if x > 0 {
				y1 := mulHeight - ((r.data[x-1]+r.data[x])*mulHeight/2)/max
				y2 := mulHeight - r.data[x]*mulHeight/max

				if y2 < y1 {
					y1, y2 = y2, y1
				}

				if uint64(y) > y1 && uint64(y) < y2 || uint64(y) == y1 {
					bitmap[y][x] = 1
				}
			}

			if x < len(row)-1 {
				y1 := mulHeight - r.data[x]*mulHeight/max
				y2 := mulHeight - ((r.data[x]+r.data[x+1])*mulHeight/2)/max

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

	for y := 0; y < int(mulHeight); y += hMult {
		line := bytes.NewBufferString("")

		for x := 0; x < int(mulWidth); x += wMult {
			bits := uint8(0)

			for by, brow := range set.bits {
				for bx, bit := range brow {
					bits = bits | (1<<bit)*bitmap[y+by][x+bx]
				}
			}

			line.WriteRune(set.runes[bits])
		}

		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}
