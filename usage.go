package boa

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// HelpFunc puts out the help for the command. Used when a user calls help [command].
// Set by calling Cobra's SetHelpFunc
func HelpFunc(cmd *cobra.Command, s []string) {
	if err := tea.NewProgram(newCmdModel(cmd), tea.WithAltScreen(), tea.WithMouseCellMotion()).Start(); err != nil {
		log.Fatal(err)
	}
}

// UsageFunc puts out the usage for the command. Used when a user provides invalid input.
// Set by calling Cobra's SetUsageFunc
func UsageFunc(cmd *cobra.Command) error {
	if err := tea.NewProgram(newCmdModel(cmd), tea.WithAltScreen(), tea.WithMouseCellMotion()).Start(); err != nil {
		return err
	}
	return nil
}
