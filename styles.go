package boa

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
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
	// default width of the help/usage output. To override supply a width to the style you want to increase
	width  = 100
	height = 50
	// Style of the border
	BorderStyle = lipgloss.NewStyle().
			Padding(0, 1, 0, 1).
			Width(width).
			BorderForeground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
			Border(lipgloss.ThickBorder())

	// Style of the title
	TitleStyle = lipgloss.NewStyle().Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
			Width(width - 4).
			Align(lipgloss.Center)
	SubTitleSytle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white}).Align(lipgloss.Center)

	// Style of the individual help sections (Exaple, Usage, Flags etc.. )
	SectionStyle = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
			Underline(true).
			BorderBottom(true).
			Margin(1, 0, 1, 0).
			Padding(0, 1, 0, 1).Align(lipgloss.Center)

	// Style of the text output
	TextStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 5).Align(lipgloss.Left).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	// Style of the selection list items
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Background(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
				Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})

	// Style of the list items
	ItemStyle = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})

	InfoStyle = lipgloss.NewStyle().Bold(true).Width(width).Align(lipgloss.Center).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})
)

func usage(cmd *cobra.Command, list list.Model) string {
	usageText := strings.Builder{}

	cmdTitle := ""
	if !cmd.HasParent() {
		rootCmdName := SectionStyle.Render(cmd.Root().Name() + " " + cmd.Root().Version)
		rootCmdLong := SubTitleSytle.Render(cmd.Root().Long)
		cmdTitle = TitleStyle.Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).
			Render(lipgloss.JoinVertical(lipgloss.Top, rootCmdName, rootCmdLong))
	}
	usageText.WriteString(cmdTitle + "\n")

	cmdSection := SectionStyle.Render("Cmd Description:")
	short := TextStyle.Render(cmd.Short)

	usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, cmdSection, short) + "\n")

	if cmd.Runnable() {
		usage := SectionStyle.Render("Usage:")
		useLine := TextStyle.Render(cmd.UseLine())
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, usage, useLine) + "\n")
		if cmd.HasAvailableSubCommands() {
			commandPath := TextStyle.Render(cmd.CommandPath() + " [command]")
			usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, commandPath) + "\n")
		}
	}

	if len(cmd.Aliases) > 0 {
		aliases := SectionStyle.Render("Aliases:")
		nameAndAlias := TextStyle.Render(cmd.NameAndAliases())
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, aliases, nameAndAlias) + "\n")
	}

	if cmd.HasAvailableLocalFlags() {
		localFlags := SectionStyle.Render("Flags:")
		flagUsage := TextStyle.Render(strings.TrimRightFunc(cmd.LocalFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, localFlags, flagUsage) + "\n")
	}
	if cmd.HasAvailableInheritedFlags() {
		globalFlags := SectionStyle.Render("Global Flags:")
		flagUsage := TextStyle.Render(strings.TrimRightFunc(cmd.InheritedFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, globalFlags, flagUsage) + "\n")
	}

	if cmd.HasExample() {
		examples := SectionStyle.Render("Examples:")
		example := TextStyle.Render(cmd.Example)
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, examples, example) + "\n")
	}

	if cmd.HasAvailableSubCommands() {
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, list.View()))
	}

	usageCard := BorderStyle.Render(usageText.String() + "\n")
	return usageCard
}

func footer(contentHeight int, windowHeight int) string {
	var help, scroll string
	help = InfoStyle.Render("↑/k up • ↓/j down • / to filter • b to go back • enter to select • q, ctrl+c to quit")
	// If content is larger than the window minus the size of the necessary footer then it will be in a scrollable viewport
	if contentHeight > windowHeight-2 {
		scroll = InfoStyle.Render("ctrl+k up • ctrl+j down • mouse to scroll")
	}
	return lipgloss.JoinVertical(lipgloss.Top, help, scroll)
}
