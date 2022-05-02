package boa

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultWidth  = 100
	defaultHeight = 150

	//default colors
	purple    = `#7e2fcc`
	darkGrey  = `#353C3B`
	lightTeal = `#03DAC5`
	darkTeal  = `#01A299`
	white     = `#e5e5e5`
	red       = `#F45353`
)

type Styles struct {
	Border       lipgloss.Style
	Title        lipgloss.Style
	SubTitle     lipgloss.Style
	Section      lipgloss.Style
	Text         lipgloss.Style
	SelectedItem lipgloss.Style
	Item         lipgloss.Style
	Info         lipgloss.Style
	CmdPrint     lipgloss.Style
}

func defaultStyles(width int) *Styles {
	s := &Styles{}

	// Style of the border
	s.Border = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Width(width).
		BorderForeground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
		Border(lipgloss.ThickBorder())

		// Style of the title
	s.Title = lipgloss.NewStyle().Bold(true).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
		Width(width - 4).
		Align(lipgloss.Center)

	// Style of the SubTitle
	s.SubTitle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white}).Align(lipgloss.Center)

	// Style of the individual help sections (Exaple, Usage, Flags etc.. )
	s.Section = lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
		Underline(true).
		BorderBottom(true).
		Margin(1, 0, 1, 0).
		Padding(0, 1, 0, 1).Align(lipgloss.Center)

	// Style of the text output
	s.Text = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 5).Align(lipgloss.Left).
		Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	// Style of the selection list items
	s.SelectedItem = lipgloss.NewStyle().PaddingLeft(2).Background(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
		Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})

	// Style of the list items
	s.Item = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})

	// Style of the info text
	s.Info = lipgloss.NewStyle().Bold(true).Width(width).Align(lipgloss.Center).
		Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	// Style of the Cmd Print text
	s.CmdPrint = lipgloss.NewStyle().Bold(true).Width(width).Margin(1).Align(lipgloss.Center).
		Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	return s
}
