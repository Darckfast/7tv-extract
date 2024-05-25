package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	PINK_COLOR = lipgloss.Color("#f629d7")
	PINK_STYLE = lipgloss.
			NewStyle().
			Bold(true).
			Foreground(PINK_COLOR)

	INFO_COLOR = lipgloss.Color("#29b6f6")
	INFO_STYLE = lipgloss.
			NewStyle().
			Bold(true)

	INFO_FG_STYLE = lipgloss.
			NewStyle().
			Bold(true).
			Foreground(INFO_COLOR)

	INFO_HG_STYLE = lipgloss.
			NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(INFO_COLOR)

	WARN_COLOR = lipgloss.Color("#ef45ab")
	WARN_STYLE = lipgloss.
			NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#fff")).
			Background(WARN_COLOR)
)

func Highlight(text string) {
	fmt.Printf("%s\n", PINK_STYLE.Render(text))
}

func Info(text string) {
	fmt.Printf("%s\n", INFO_STYLE.Render(text))
}

func Warn(text string) {
	fmt.Printf("%s\n", WARN_STYLE.Render(text))
}

func Progress(curr, total int) {
	fd := int(os.Stdout.Fd())
	termWidth, _, _ := term.GetSize(fd)
	numProg := fmt.Sprintf("[%d/%d] ", curr, total)
	charTotal := termWidth - len(numProg)
	prog := int((float32(curr) / float32(total)) * float32(charTotal))

	if charTotal < 0 {
		charTotal = 0
	}

	spacesChars := strings.Repeat(" ", charTotal-prog)
	progChars := strings.Repeat(INFO_FG_STYLE.Render("░"), prog)

	fmt.Printf("\r%s%s%s", numProg, progChars, spacesChars)
}
