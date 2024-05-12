package utils

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func PrintLine(text string) {
	fd := int(os.Stdout.Fd())
	termWidth, _, _ := term.GetSize(fd)
	charTotal := termWidth - len(text) - 1

	if charTotal < 0 {
		charTotal = 0
	}

	fmt.Printf("\r%s %s", text, strings.Repeat(" ", charTotal))
}
