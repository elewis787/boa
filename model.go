package boa

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const spacebar = " "

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
	vp := viewport.New(0, 0)
	vp.KeyMap = viewPortKeyMap()
	return &cmdModel{
		cmd:           cmd,
		subCmds:       subCmds,
		list:          l,
		viewport:      &vp,
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
		case "b":
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
	m.viewport.KeyMap = viewPortKeyMap()
	if m.viewport.Height > m.contentHeight+lipgloss.Height(footer(m.contentHeight, m.windowHeight)) {
		m.viewport.Height = m.contentHeight
	}
	cmds = append(cmds, listCmd, viewPortCmd)
	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string.
func (m cmdModel) View() string {
	m.viewport.SetContent(usage(m.cmd, m.list))
	return lipgloss.JoinVertical(lipgloss.Top, m.viewport.View(), footer(m.contentHeight, m.windowHeight))
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
