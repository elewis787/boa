package boa

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	purple    = `#7e2fcc`
	darkGrey  = `#353C3B`
	lightTeal = `#03DAC5`
	darkTeal  = `#01A299`
	white     = `#e5e5e5`
	red       = `#F45353`
)

var (
	// TODO func to calc width
	physicalWidth, _, _ = term.GetSize(int(os.Stdout.Fd()))
	width               = physicalWidth / 3

	titleStyle = lipgloss.NewStyle().Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
			Padding(1, 1).
			Align(lipgloss.Center)

	sectionStyle = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
			Underline(true).
			BorderBottom(true).
			Padding(0, 1, 0, 1).Align(lipgloss.Center)

	textStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 5).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	subTextStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 2).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})
)
