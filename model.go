package boa

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const listHeight = 10

type cmdItem string

type item struct {
	cmd *cobra.Command
}

func (i item) FilterValue() string { return i.cmd.Name() }

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := i.cmd.Name() + lipgloss.NewStyle().Bold(true).
		PaddingLeft(i.cmd.NamePadding()-len(i.cmd.Name())+1).Render(i.cmd.Short)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}
	fmt.Fprintf(w, fn(str))
}

type itemDelegate struct{}

type CmdModel struct {
	list    list.Model
	cmd     *cobra.Command
	subCmds []list.Item
}

func getSubCommands(c *cobra.Command) []list.Item {
	subs := make([]list.Item, 0)
	if c.HasAvailableSubCommands() {
		for _, subcmd := range c.Commands() {
			if subcmd.Name() == "help" || subcmd.IsAvailableCommand() {
				subs = append(subs, item{cmd: subcmd})
			}
		}
	}
	return subs
}

func newList(items []list.Item) list.Model {
	l := list.New(items, itemDelegate{}, 0, listHeight)
	l.Styles.TitleBar.Padding(0, 0)
	l.Styles.Title = sectionStyle
	l.Title = "Available Sub Commands:"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	return l
}

func NewCmdModel(cmd *cobra.Command) *CmdModel {
	subCmds := getSubCommands(cmd)
	l := newList(subCmds)
	return &CmdModel{cmd: cmd, subCmds: subCmds, list: l}
}

func (m CmdModel) Init() tea.Cmd {
	return nil
}

func (m CmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.list = newList(subCmds)
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m CmdModel) View() string {
	cmdTitle := ""
	if !m.cmd.HasParent() {
		rootCmdName := sectionStyle.Render(m.cmd.Root().Name() + " " + m.cmd.Root().Version)
		rootCmdLong := lipgloss.NewStyle().Align(lipgloss.Center).Render(m.cmd.Root().Long)

		cmdTitle = titleStyle.Width(100).Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).Render(lipgloss.JoinVertical(lipgloss.Top, rootCmdName, rootCmdLong))
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

	if m.cmd.HasAvailableSubCommands() {
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, m.list.View())
		subCmdHelp := subTextStyle.Render("\nUse \"" + m.cmd.CommandPath() + " [command] --help\" for more information about a command.")
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, subCmdHelp)
	}

	usageOutput = lipgloss.JoinVertical(lipgloss.Top, cmdTitle, usageOutput)
	usageOutput = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		BorderForeground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
		Border(lipgloss.ThickBorder()).Render(usageOutput)

	closeText := textStyle.Render("Press q, ctrl+c or esc to close")
	return lipgloss.JoinVertical(lipgloss.Top, usageOutput, closeText, "\n")

}
