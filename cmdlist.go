package boa

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// screen length of the list
const listHeight = 10

// item is the object that will appear in our list
type item struct {
	cmd *cobra.Command
}

func (i item) FilterValue() string { return i.cmd.Name() }

// itemDelegate encapsulates the general functionality for all list items
type itemDelegate struct{}

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

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return SelectedItemStyle.Render("> " + s)
		}
	}
	fmt.Fprintf(w, fn(str))
}

// newSubCmdsList returns a new list.Model filled with the values in []list.Items
func newSubCmdsList(items []list.Item) list.Model {
	l := list.New(items, itemDelegate{}, 0, listHeight)
	l.Styles.TitleBar.Padding(0, 0)
	l.Styles.Title = SectionStyle
	l.Title = "Available Sub Commands:"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	return l
}

// getSubCommands returns a []list.Item filled with any available sub command from the supplied *cobra.Command.
// This does not follow the command chain past a depth of 1.
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
