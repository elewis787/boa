package boa

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// Inital options
var defaultOpts = defaultOptions()

// HelpFunc puts out the help for the command. Used when a user calls help [command].
// Set by calling Cobra's SetHelpFunc
func HelpFunc(cmd *cobra.Command, s []string) {
	model := newCmdModel(defaultOpts, defaultStyles(defaultOpts.width), cmd)
	if err := tea.NewProgram(model, defaultOpts.atlScreen, defaultOpts.mouseCellMotion).Start(); err != nil {
		log.Fatal(err)
	}
	if model.print {
		fmt.Println(model.cmdChain)
	}
}

// UsageFunc puts out the usage for the command. Used when a user provides invalid input.
// Set by calling Cobra's SetUsageFunc
func UsageFunc(cmd *cobra.Command) error {
	model := newCmdModel(defaultOpts, defaultStyles(defaultOpts.width), cmd)
	if err := tea.NewProgram(model, defaultOpts.atlScreen, defaultOpts.mouseCellMotion).Start(); err != nil {
		return err
	}
	if model.print {
		fmt.Println(model.cmdChain)
	}
	return nil
}
