package boa

import (
	"github.com/charmbracelet/lipgloss"
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
	width = 100

	titleStyle = lipgloss.NewStyle().Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
			Padding(1, 1).
			Align(lipgloss.Center)

	sectionStyle = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
			Underline(true).
			BorderBottom(true).
			Margin(0, 0, 1, 0).
			Padding(0, 1, 0, 1).Align(lipgloss.Center)

	textStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 5).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	subTextStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 2).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Background(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
				Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})
	itemStyle = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})
)
