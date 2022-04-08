package boa

import (
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type CmdModel struct {
	cmd *cobra.Command
}

func NewCmdModel(cmd *cobra.Command) *CmdModel {
	return &CmdModel{cmd}
}

func (m CmdModel) Init() tea.Cmd {
	return cmdFunc(m.cmd)
}

func (m CmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case *cobra.Command:
		return m, cmdFunc(msg)
	}
	return m, nil
}

func (m CmdModel) View() string {
	cmdTitle := ""
	if !m.cmd.HasParent() {
		rootCmdName := sectionStyle.Render(m.cmd.Root().Name() + " " + m.cmd.Root().Version)
		rootCmdLong := lipgloss.NewStyle().Align(lipgloss.Center).Render(m.cmd.Root().Long)

		cmdTitle = titleStyle.Width(width).Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).Render(lipgloss.JoinVertical(lipgloss.Top, rootCmdName, rootCmdLong))
	}

	usageOutput := sectionStyle.Render("Cmd Description:")
	short := textStyle.Render(m.cmd.Short)

	usage := sectionStyle.Render("Usage:")
	usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, short, usage)

	if m.cmd.Runnable() {
		useLine := textStyle.Render(m.cmd.UseLine())
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, useLine)
	}

	if m.cmd.HasAvailableSubCommands() {
		commandPath := textStyle.Render(m.cmd.CommandPath() + " [command]")
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, commandPath)
	}

	if len(m.cmd.Aliases) > 0 {
		aliases := sectionStyle.Render("Aliases:")
		nameAndAlias := textStyle.Render(m.cmd.NameAndAliases())
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, aliases, nameAndAlias)

	}

	if m.cmd.HasAvailableSubCommands() {
		commands := sectionStyle.Render("Available Commands:")
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, commands)

		for _, subcmd := range m.cmd.Commands() {
			if subcmd.Name() == "help" || subcmd.IsAvailableCommand() {
				cmd := textStyle.Render(subcmd.Name()) +
					lipgloss.NewStyle().
						Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).Bold(true).
						PaddingLeft(subcmd.NamePadding()-len(subcmd.Name())+1).Render(subcmd.Short)
				usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, cmd)
			}
		}
		subCmdHelp := subTextStyle.Render("\nUse \"" + m.cmd.CommandPath() + " [command] --help\" for more information about a command.")
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, subCmdHelp)
	}

	if m.cmd.HasAvailableLocalFlags() {
		localFlags := sectionStyle.Render("Flags:")
		flagUsage := textStyle.Render(strings.TrimFunc(m.cmd.LocalFlags().FlagUsages(), unicode.IsSpace))
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, localFlags, flagUsage)
	}

	if m.cmd.HasAvailableInheritedFlags() {
		globalFlags := sectionStyle.Render("Global Flags:")
		flagUsage := textStyle.Render(strings.TrimFunc(m.cmd.InheritedFlags().FlagUsages(), unicode.IsSpace))
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, globalFlags, flagUsage)
	}

	if m.cmd.HasExample() {
		examples := sectionStyle.Render("Examples:")
		example := textStyle.Render(m.cmd.Example)
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, examples, example)
	}
	usageOutput = lipgloss.JoinVertical(lipgloss.Top, cmdTitle, usageOutput)
	usageOutput = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		BorderForeground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
		Border(lipgloss.ThickBorder()).Render(usageOutput)

	return usageOutput
}

func cmdFunc(cmd *cobra.Command) tea.Cmd {
	return func() tea.Msg {
		return cmd
	}
}
