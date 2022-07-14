package boa

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	spacebar = " "
)

// cmdModel Implements tea.Model. It provides an interactive
// help and usage tui component for bubbletea programs.
type cmdModel struct {
	styles   *Styles
	list     list.Model
	viewport *viewport.Model
	cmd      *cobra.Command
	subCmds  []list.Item
	print    bool
	cmdChain string
	// Store window height to adjust viewport on command selection changes
	windowHeight int
	// Store full height of content for given view, updated on command change
	contentHeight int
	errorWriter   *bytes.Buffer
}

// newCmdModel initializes a based on values supplied from cmd *cobra.Command
func newCmdModel(options *options, cmd *cobra.Command) *cmdModel {
	subCmds := getSubCommands(cmd)
	l := newSubCmdsList(options.styles, subCmds)
	vp := viewport.New(0, 0)
	vp.KeyMap = viewPortKeyMap()
	m := &cmdModel{
		styles:      options.styles,
		cmd:         cmd,
		subCmds:     subCmds,
		list:        l,
		viewport:    &vp,
		errorWriter: options.errorWriter,
	}
	m.contentHeight = lipgloss.Height(m.usage())
	return m
}

// Init is the initial cmd to be executed which is nil for this component.
func (m *cmdModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received.
func (m *cmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var listCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = m.windowHeight - lipgloss.Height(m.footer())
		if m.viewport.Height > m.contentHeight {
			m.viewport.Height = m.contentHeight
		}
		// Scroll viewport back to top for new screen
		m.viewport.SetYOffset(0)
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
				m.list = newSubCmdsList(m.styles, subCmds)
				m.viewport.Height = m.windowHeight - lipgloss.Height(m.footer())
				// Update new content height and check viewport size
				m.contentHeight = lipgloss.Height(m.usage())
				if m.viewport.Height > m.contentHeight {
					m.viewport.Height = m.contentHeight
				}
				// Scroll viewport back to top for new screen
				m.viewport.SetYOffset(0)
			}
			return m, nil
		case "b":
			if m.cmd.HasParent() {
				m.cmd = m.cmd.Parent()
				subCmds := getSubCommands(m.cmd)
				m.list = newSubCmdsList(m.styles, subCmds)
				m.viewport.Height = m.windowHeight - lipgloss.Height(m.footer())
				// Update new content height and check viewport size
				m.contentHeight = lipgloss.Height(m.usage())
				if m.viewport.Height > m.contentHeight {
					m.viewport.Height = m.contentHeight
				}
				// Scroll viewport back to top for new screen
				m.viewport.SetYOffset(0)
			}
			return m, nil
		case "p":
			m.print = true
			m.cmdChain = print("", m.cmd)
			return m, tea.Quit
		}
	}

	m.list, listCmd = m.list.Update(msg)
	newViewport, viewPortCmd := m.viewport.Update(msg)
	// point to new viewport
	m.viewport = &newViewport
	m.viewport.KeyMap = viewPortKeyMap()
	if m.viewport.Height > m.contentHeight {
		m.viewport.Height = m.contentHeight
	}
	cmds = append(cmds, listCmd, viewPortCmd)
	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string.
func (m *cmdModel) View() string {
	m.viewport.SetContent(m.usage())
	return lipgloss.JoinVertical(lipgloss.Top, m.viewport.View(), m.footer())
}

// usage builds the usage body from a cobra command
func (m *cmdModel) usage() string {
	usageText := strings.Builder{}

	if m.errorWriter != nil && m.errorWriter.Len() > 0 {
		usageText.WriteString(m.styles.ErrorText.Render(m.errorWriter.String() + "\n"))
	}

	cmdTitle := ""
	cmdName := m.cmd.Name()
	if m.cmd.Version != "" {
		cmdName += " " + m.cmd.Version
	}
	cmdName = m.styles.Section.Render(cmdName)
	cmdLong := m.styles.SubTitle.Render(m.cmd.Long)
	cmdTitle = m.styles.Title.Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).
		Render(lipgloss.JoinVertical(lipgloss.Center, cmdName, cmdLong))
	usageText.WriteString(cmdTitle + "\n")

	cmdSection := m.styles.Section.Render("Cmd Description:")
	short := m.styles.Text.Render(m.cmd.Short)

	usageText.WriteString(lipgloss.JoinVertical(lipgloss.Left, cmdSection, short) + "\n")

	if m.cmd.Runnable() {
		usage := m.styles.Section.Render("Usage:")
		useLine := m.styles.Text.Render(m.cmd.UseLine())
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Left, usage, useLine) + "\n")
		if m.cmd.HasAvailableSubCommands() {
			commandPath := m.styles.Text.Render(m.cmd.CommandPath() + " [command]")
			usageText.WriteString(lipgloss.JoinVertical(lipgloss.Left, commandPath) + "\n")
		}
	}

	if len(m.cmd.Aliases) > 0 {
		aliases := m.styles.Section.Render("Aliases:")
		nameAndAlias := m.styles.Text.Render(m.cmd.NameAndAliases())
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Left, aliases, nameAndAlias) + "\n")
	}

	if m.cmd.HasAvailableLocalFlags() {
		localFlags := m.styles.Section.Render("Flags:")
		flagUsage := m.styles.Text.Render(strings.TrimRightFunc(m.cmd.LocalFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Left, localFlags, flagUsage) + "\n")
	}
	if m.cmd.HasAvailableInheritedFlags() {
		globalFlags := m.styles.Section.Render("Global Flags:")
		flagUsage := m.styles.Text.Render(strings.TrimRightFunc(m.cmd.InheritedFlags().FlagUsages(), unicode.IsSpace))
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Left, globalFlags, flagUsage) + "\n")
	}

	if m.cmd.HasExample() {
		examples := m.styles.Section.Render("Examples:")
		example := m.styles.Text.Render(m.cmd.Example)
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Left, examples, example) + "\n")
	}

	if m.cmd.HasAvailableSubCommands() {
		usageText.WriteString(lipgloss.JoinVertical(lipgloss.Left, m.list.View()))
	}

	return m.styles.Border.Render(usageText.String() + "\n")
}

// footer outputs the footer of the viewport and contains help text.
func (m *cmdModel) footer() string {
	var help, scroll string
	help = m.styles.Info.Render("↑/k up • ↓/j down • / to filter • p to print • b to go back • enter to select • q, ctrl+c to quit")
	// If content is larger than the window minus the size of the necessary footer then it will be in a scrollable viewport
	if m.contentHeight > m.windowHeight-2 {
		scroll = m.styles.Info.Render("ctrl+k up • ctrl+j down • mouse to scroll")
	}
	return lipgloss.JoinVertical(lipgloss.Left, help, scroll)
}

// print outputs the command chain for a given cobra command.
func print(v string, cmd *cobra.Command) string {
	if cmd != nil {
		v = cmd.Name() + " " + v
		if !cmd.HasParent() {
			// final result
			return "Command: " + v
		}
		// recursively walk cmd chain
		return print(v, cmd.Parent())
	}
	return v
}

func viewPortKeyMap() viewport.KeyMap {
	return viewport.KeyMap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", spacebar, "f"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "v"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
		),
		Up: key.NewBinding(
			key.WithKeys("ctrl+k"),
		),
		Down: key.NewBinding(
			key.WithKeys("ctrl+j"),
		),
	}
}
