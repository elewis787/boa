package boa

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func HelpFunc(cmd *cobra.Command, s []string) {
	if err := tea.NewProgram(newCmdModel(cmd)).Start(); err != nil {
		log.Fatal(err)
	}
}

func UsageFunc(cmd *cobra.Command) error {
	if err := tea.NewProgram(newCmdModel(cmd)).Start(); err != nil {
		return err
	}
	return nil
}
