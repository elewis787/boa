package boa

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/cobra"
)

// cmdModel Implements tea.Model. It provides an interactive
// help and usage tui component for bubbletea programs.
type cmdModel struct {
	list     list.Model
	ready    bool
	viewport viewport.Model
	cmd      *cobra.Command
	subCmds  []list.Item
	cursor   int
}

// newCmdModel initializes a based on values supplied from cmd *cobra.Command
func newCmdModel(cmd *cobra.Command) *cmdModel {
	subCmds := getSubCommands(cmd)
	l := newSubCmdsList(subCmds)
	l.KeyMap.CursorDown = key.NewBinding(key.WithKeys("m"))
	l.KeyMap.CursorUp = key.NewBinding(key.WithKeys("n"))
	return &cmdModel{
		cmd:     cmd,
		subCmds: subCmds,
		list:    l,
	}
}

// Init is the initial cmd to be executed which is nil for this component.
func (m cmdModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received.
func (m cmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var viewPortCmd, listCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.Model{
				Width:  msg.Width,
				Height: msg.Height - 2,
			}
			m.viewport.SetContent(m.p())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 2
		}
		return m, nil
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
				m.viewport.SetContent(m.p())
			}
			return m, nil
		}
	}
	m.list, listCmd = m.list.Update(msg)
	m.viewport, viewPortCmd = m.viewport.Update(msg)
	cmds = append(cmds, listCmd, viewPortCmd)
	m.viewport.SetContent(m.p())

	return m, tea.Batch(cmds...)
}

func (m cmdModel) View() string {
	return m.viewport.View()
}

// View renders the program's UI, which is just a string.
func (m cmdModel) p() string {
	usageText := strings.Builder{}

	cmdTitle := ""
	if !m.cmd.HasParent() {
		rootCmdName := SectionStyle.Render(m.cmd.Root().Name() + " " + m.cmd.Root().Version)
		rootCmdLong := lipgloss.NewStyle().Align(lipgloss.Center).Render(m.cmd.Root().Long)
		cmdTitle = TitleStyle.Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).
			Render(lipgloss.JoinVertical(lipgloss.Top, rootCmdName, rootCmdLong))
	}
	usageText.WriteString(cmdTitle + "\n")

	cmdSection := SectionStyle.Render("Cmd Description:")
	short := TextStyle.Render(m.cmd.Short)

	usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, cmdSection, short) + "\n")

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
		flagUsage := TextStyle.Render(strings.TrimRightFunc(m.cmd.LocalFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, localFlags, flagUsage) + "\n")
	}
	if m.cmd.HasAvailableLocalFlags() {
		localFlags := SectionStyle.Render("Flags:")
		flagUsage := TextStyle.Render(strings.TrimRightFunc(m.cmd.LocalFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, localFlags, flagUsage) + "\n")
	}
	if m.cmd.HasAvailableInheritedFlags() {
		globalFlags := SectionStyle.Render("Global Flags:")
		flagUsage := TextStyle.Render(strings.TrimRightFunc(m.cmd.InheritedFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, globalFlags, flagUsage) + "\n")
	}

	if m.cmd.HasExample() {
		examples := SectionStyle.Render("Examples:")
		example := TextStyle.Render(m.cmd.Example)
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, examples, example) + "\n")
	}

	if m.cmd.HasAvailableSubCommands() {
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Top, m.list.View()))
	}

	usageCard := BorderStyle.Render(usageText.String() + "\n")
	return lipgloss.JoinVertical(lipgloss.Top, usageCard, InfoStyle.Render("↑/k up • ↓/j down • / to filter • enter to select • q, ctrl+c, esc to quit"), "\n\n")
}

func formatFlags(flagText string) string {
	text := ""
	for _, line := range strings.Split(strings.TrimRight(flagText, "\n"), "\n") {
		if len(line) > width {
			tmp := wordwrap.String(line, width) + "\n"
			for i, v := range strings.Split(strings.TrimRight(tmp, "\n"), "\n") {
				if i > 0 {
					v = indent.String(v, 6)
				}
				text += v + "\n"
			}
		} else {
			text += line + "\n"
		}
	}
	return text
}
