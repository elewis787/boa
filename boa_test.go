package boa

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/spf13/cobra"
)

func TestNew(t *testing.T) {
	// Default instance
	b := New()
	if b == nil {
		t.Fatal("expected non-nil Boa instance")
	}
	if b.options == nil {
		t.Fatal("expected non-nil options")
	}

	// Nil option should not panic
	bWithNil := New(nil)
	if bWithNil == nil {
		t.Fatal("expected non-nil Boa instance when nil option provided")
	}

	// Custom options
	customBuf := &bytes.Buffer{}
	customStyles := DefaultStyles()
	bCustom := New(
		WithAltScreen(true),
		WithStyles(customStyles),
		WithStyles(nil), // nil styles should preserve existing styles
		WithErrWriter(customBuf),
		WithErrWriter(nil), // nil writer should preserve existing writer
	)
	if bCustom.options.errorWriter != customBuf {
		t.Errorf("expected custom errorWriter to be set")
	}
}

func TestGetSubCommands(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	subCmd1 := &cobra.Command{Use: "sub1", Short: "First subcommand", Run: func(cmd *cobra.Command, args []string) {}}
	subCmd2 := &cobra.Command{Use: "sub2", Short: "Second subcommand", Run: func(cmd *cobra.Command, args []string) {}}
	hiddenCmd := &cobra.Command{Use: "hidden", Hidden: true, Run: func(cmd *cobra.Command, args []string) {}}

	rootCmd.AddCommand(subCmd1, subCmd2, hiddenCmd)

	items := getSubCommands(rootCmd)
	if len(items) != 2 {
		t.Fatalf("expected 2 visible subcommands, got %d", len(items))
	}

	// Leaf command without subcommands
	emptyItems := getSubCommands(subCmd1)
	if len(emptyItems) != 0 {
		t.Fatalf("expected 0 subcommands for leaf command, got %d", len(emptyItems))
	}
}

func TestItemFilterValue(t *testing.T) {
	cmd := &cobra.Command{Use: "my-cmd"}
	i := item{cmd: cmd}
	if i.FilterValue() != "my-cmd" {
		t.Errorf("expected filter value 'my-cmd', got %q", i.FilterValue())
	}
}

func TestItemDelegateRender(t *testing.T) {
	styles := DefaultStyles()
	cmd := &cobra.Command{Use: "testcmd", Short: "Test description"}
	it := item{cmd: cmd}
	delegate := itemDelegate{styles: styles}

	if delegate.Height() != 1 {
		t.Errorf("expected Height() 1, got %d", delegate.Height())
	}
	if delegate.Spacing() != 0 {
		t.Errorf("expected Spacing() 0, got %d", delegate.Spacing())
	}

	// Render unselected item
	items := []list.Item{it}
	l := newSubCmdsList(styles, items)

	bufUnselected := &bytes.Buffer{}
	delegate.Render(bufUnselected, l, 1, it) // index 1 != selected index 0
	outUnselected := bufUnselected.String()
	if !strings.Contains(outUnselected, "testcmd") {
		t.Errorf("expected unselected output to contain command name, got %q", outUnselected)
	}

	// Render selected item (must NOT have bracket "[testcmd...]")
	bufSelected := &bytes.Buffer{}
	delegate.Render(bufSelected, l, 0, it) // index 0 == selected index 0
	outSelected := bufSelected.String()
	if !strings.Contains(outSelected, "> ") {
		t.Errorf("expected selected output to contain prefix '> ', got %q", outSelected)
	}
	if strings.Contains(outSelected, "[") || strings.Contains(outSelected, "]") {
		t.Errorf("selected item output should not contain brackets, got %q", outSelected)
	}
}

func TestDefaultStyles(t *testing.T) {
	styles := DefaultStyles()
	if styles == nil {
		t.Fatal("expected non-nil default styles")
	}
}
