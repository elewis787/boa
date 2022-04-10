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
	Width = 100

	TitleStyle = lipgloss.NewStyle().Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
			Width(Width - 4).
			Align(lipgloss.Center)

	SectionStyle = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
			Underline(true).
			BorderBottom(true).
			Margin(1, 0, 1, 0).
			Padding(0, 1, 0, 1).Align(lipgloss.Center)

	TextStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 5).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Background(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
				Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})
	ItemStyle = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})

	InfoStyle = lipgloss.NewStyle().Bold(true).Width(Width).Align(lipgloss.Center).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})
)
