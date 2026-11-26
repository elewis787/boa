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
	baseStyles *Styles
	styles     *Styles
	list       list.Model
	viewport   *viewport.Model
	cmd        *cobra.Command
	subCmds    []list.Item
	print      bool
	cmdChain   string
	// Store window dimensions to adjust viewport and responsive styles
	windowWidth  int
	windowHeight int
	// Store full height of content for given view, updated on command change
	contentHeight int
	errorWriter   *bytes.Buffer
}

// newCmdModel initializes a based on values supplied from cmd *cobra.Command
func newCmdModel(options *options, cmd *cobra.Command) *cmdModel {
	if options == nil {
		options = defaultOptions()
	}
	if options.styles == nil {
		options.styles = DefaultStyles()
	}
	subCmds := getSubCommands(cmd)
	l := newSubCmdsList(options.styles, subCmds)
	vp := viewport.New(0, 0)
	vp.KeyMap = viewPortKeyMap()
	m := &cmdModel{
		baseStyles:  options.styles,
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
		m.windowWidth = msg.Width
		m.styles = fitStylesToWidth(m.baseStyles, msg.Width)
		if borderW := m.styles.Border.GetWidth(); borderW > 2 {
			m.list.SetWidth(borderW - 2)
		}
		m.contentHeight = lipgloss.Height(m.usage())
		m.viewport.Width = msg.Width
		m.viewport.Height = max(0, m.windowHeight-lipgloss.Height(m.footer()))
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
				m.refreshLayout(getSubCommands(m.cmd))
			}
			return m, nil
		case "b":
			if m.cmd.HasParent() {
				m.cmd = m.cmd.Parent()
				m.refreshLayout(getSubCommands(m.cmd))
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
	if m == nil || m.cmd == nil {
		return ""
	}

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

	subTitleStyle := m.styles.SubTitle
	if subTitleStyle.GetWidth() == 0 && m.styles.Title.GetWidth() > 2 {
		subTitleStyle = subTitleStyle.Width(m.styles.Title.GetWidth() - 2)
	}
	cmdLong := subTitleStyle.Render(m.cmd.Long)
	cmdTitle = m.styles.Title.Render(lipgloss.JoinVertical(lipgloss.Center, cmdName, cmdLong))
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

func (m *cmdModel) refreshLayout(subCmds []list.Item) {
	m.list = newSubCmdsList(m.styles, subCmds)
	if m.windowWidth > 0 {
		if borderW := m.styles.Border.GetWidth(); borderW > 2 {
			m.list.SetWidth(borderW - 2)
		}
	}
	m.contentHeight = lipgloss.Height(m.usage())
	m.viewport.Height = max(0, m.windowHeight-lipgloss.Height(m.footer()))
	if m.viewport.Height > m.contentHeight {
		m.viewport.Height = m.contentHeight
	}
	m.viewport.SetYOffset(0)
}

func fitStylesToWidth(base *Styles, windowWidth int) *Styles {
	if base == nil {
		base = DefaultStyles()
	}
	if windowWidth <= 0 {
		return base
	}

	baseBorderWidth := base.Border.GetWidth()
	if baseBorderWidth <= 0 {
		baseBorderWidth = defaultWidth
	}

	maxInner := windowWidth - 2
	if maxInner >= baseBorderWidth && (base.Info.GetWidth() <= 0 || base.Info.GetWidth() <= windowWidth) {
		return base
	}

	targetWidth := min(baseBorderWidth, maxInner)
	if targetWidth < 10 {
		targetWidth = 10
	}

	// Copy struct so base (options.styles) is never mutated
	s := *base
	s.Border = base.Border.Width(targetWidth)
	s.Title = base.Title.Width(max(1, targetWidth-4))
	s.SubTitle = base.SubTitle.Width(max(1, targetWidth-6))
	s.ErrorText = base.ErrorText.Width(max(1, targetWidth-4))
	s.CmdPrint = base.CmdPrint.Width(targetWidth)
	if infoW := base.Info.GetWidth(); infoW == 0 || infoW > windowWidth {
		s.Info = base.Info.Width(windowWidth)
	}
	return &s
}
