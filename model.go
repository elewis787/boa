package boa

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// cmdModel Implements tea.Model. It provides an interactive
// help and usage tui component for bubbletea programs.
type cmdModel struct {
	list    list.Model
	cmd     *cobra.Command
	subCmds []list.Item
}

// newCmdModel initializes a based on values supplied from cmd *cobra.Command
func newCmdModel(cmd *cobra.Command) *cmdModel {
	subCmds := getSubCommands(cmd)
	l := newSubCmdsList(subCmds)
	return &cmdModel{cmd: cmd, subCmds: subCmds, list: l}
}

// Init is the inital cmd to be executed which is nil for this component.
func (m cmdModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received.
func (m cmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.cmd = i.cmd
				subCmds := getSubCommands(i.cmd)
				m.list = newSubCmdsList(subCmds)
			}
			return m, nil
		}
	}
	// default behavior is to return our list model
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the program's UI, which is just a string.
func (m cmdModel) View() string {
	usageText := strings.Builder{}

	cmdTitle := ""
	if !m.cmd.HasParent() {
		rootCmdName := SectionStyle.Render(m.cmd.Root().Name() + " " + m.cmd.Root().Version)
		rootCmdLong := lipgloss.NewStyle().Align(lipgloss.Center).Render(m.cmd.Root().Long)
		cmdTitle = TitleStyle.Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).
			Render(lipgloss.JoinVertical(lipgloss.Top, rootCmdName, rootCmdLong))
	}

	cmdSection := SectionStyle.Render("Cmd Description:")
	short := TextStyle.Render(m.cmd.Short)

	usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, cmdTitle, cmdSection, short) + "\n")

	if m.cmd.Runnable() {
		usage := SectionStyle.Render("Usage:")
		useLine := TextStyle.Render(m.cmd.UseLine())
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, usage, useLine) + "\n")
		if m.cmd.HasAvailableSubCommands() {
			commandPath := TextStyle.Render(m.cmd.CommandPath() + " [command]")
			usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, commandPath) + "\n")
		}
	}

	if len(m.cmd.Aliases) > 0 {
		aliases := SectionStyle.Render("Aliases:")
		nameAndAlias := TextStyle.Render(m.cmd.NameAndAliases())
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, aliases, nameAndAlias) + "\n")

	}

	if m.cmd.HasAvailableLocalFlags() {
		localFlags := SectionStyle.Render("Flags:")
		flagUsage := TextStyle.Render(strings.TrimFunc(m.cmd.LocalFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, localFlags, flagUsage) + "\n")
	}

	if m.cmd.HasAvailableInheritedFlags() {
		globalFlags := SectionStyle.Render("Global Flags:")
		flagUsage := TextStyle.Render(strings.TrimFunc(m.cmd.InheritedFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, globalFlags, flagUsage) + "\n")
	}

	if m.cmd.HasExample() {
		examples := SectionStyle.Render("Examples:")
		example := TextStyle.Render(m.cmd.Example)
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, examples, example) + "\n")
	}

	if m.cmd.HasAvailableSubCommands() {
		usageText.WriteString(m.list.View() + "\n")
	}

	usageCard := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Width(Width).
		BorderForeground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
		Border(lipgloss.ThickBorder()).Render(usageText.String() + "\n")

	return lipgloss.JoinVertical(lipgloss.Top, usageCard, InfoStyle.Render("Press q, ctrl+c or esc to close"), "\n\n")
}
