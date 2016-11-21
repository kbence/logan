package utils

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func GetTerminalDimensions() (int, int) {
	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		width = 80
		height = 25
	}

	return width, height
}
