package boa

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// cmdModel Implements tea.Model. It provides an interactive
// help and usage tui component for bubbletea programs.
type cmdModel struct {
	list     list.Model
	viewport *viewport.Model
	cmd      *cobra.Command
	subCmds  []list.Item
	// Store window height to adjust viewport on command selection changes
	windowHeight int
	// Store full height of content for given view, updated on command change
	contentHeight int
}

// newCmdModel initializes a based on values supplied from cmd *cobra.Command
func newCmdModel(cmd *cobra.Command) *cmdModel {
	subCmds := getSubCommands(cmd)
	l := newSubCmdsList(subCmds)
	return &cmdModel{
		cmd:           cmd,
		subCmds:       subCmds,
		list:          l,
		viewport:      &viewport.Model{},
		contentHeight: lipgloss.Height(usage(cmd, l)),
	}
}

// Init is the initial cmd to be executed which is nil for this component.
func (m cmdModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received.
func (m cmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var listCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = m.windowHeight - lipgloss.Height(footer(m.contentHeight, m.windowHeight))
		if m.viewport.Height > m.contentHeight+lipgloss.Height(footer(m.contentHeight, m.windowHeight)) {
			m.viewport.Height = m.contentHeight
		}
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.cmd = i.cmd
				subCmds := getSubCommands(m.cmd)
				m.list = newSubCmdsList(subCmds)
				m.viewport.Height = m.windowHeight - lipgloss.Height(footer(m.contentHeight, m.windowHeight))
				// Update new content height and check viewport size
				m.contentHeight = lipgloss.Height(usage(m.cmd, m.list))
				if m.viewport.Height > m.contentHeight+lipgloss.Height(footer(m.contentHeight, m.windowHeight)) {
					m.viewport.Height = m.contentHeight
				}
				// Scroll viewport back to top for new screen
				m.viewport.SetYOffset(0)
			}
			return m, nil
		case "backspace":
			if m.cmd.HasParent() {
				m.cmd = m.cmd.Parent()
				subCmds := getSubCommands(m.cmd)
				m.list = newSubCmdsList(subCmds)
				m.viewport.Height = m.windowHeight - lipgloss.Height(footer(m.contentHeight, m.windowHeight))
				// Update new content height and check viewport size
				m.contentHeight = lipgloss.Height(usage(m.cmd, m.list))
				if m.viewport.Height > m.contentHeight+lipgloss.Height(footer(m.contentHeight, m.windowHeight)) {
					m.viewport.Height = m.contentHeight
				}
				// Scroll viewport back to top for new screen
				m.viewport.SetYOffset(0)
			}
			return m, nil
		}
	}

	m.list, listCmd = m.list.Update(msg)
	newViewport, viewPortCmd := m.viewport.Update(msg)
	// point to new viewport
	m.viewport = &newViewport
	if m.viewport.Height > m.contentHeight+lipgloss.Height(footer(m.contentHeight, m.windowHeight)) {
		m.viewport.Height = m.contentHeight
	}
	cmds = append(cmds, listCmd, viewPortCmd)
	return m, tea.Batch(cmds...)
}

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

// View renders the program's UI, which is just a string.
func (m cmdModel) View() string {
	m.viewport.SetContent(usage(m.cmd, m.list))
	return lipgloss.JoinVertical(lipgloss.Top, m.viewport.View(), footer(m.contentHeight, m.windowHeight))
}

func footer(contentHeight int, windowHeight int) string {
	help := InfoStyle.Render("↑/k up • ↓/j down • / to filter • backspace to go back • enter to select • q, ctrl+c to quit")
	scroll := InfoStyle.Render("use mouse scroll or space to see full usage")
	return lipgloss.JoinVertical(lipgloss.Top, help, scroll)
}
