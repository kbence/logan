package utils

import terminal "github.com/wayneashleyberry/terminal-dimensions"

func GetTerminalDimensions() (int, int) {
	width, err := terminal.Width()
	if err != nil {
		width = 80
	}

	height, err := terminal.Height()
	if err != nil {
		height = 25
	}

	return int(width), int(height)
}
