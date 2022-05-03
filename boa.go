package boa

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type Boa struct {
	options *options
}

// Create a new instance of Boa with custom Options
func New(options ...Options) *Boa {
	opts := defaultOptions()
	for _, opt := range options {
		opt.apply(opts)
	}
	return &Boa{
		options: opts,
	}
}

func (b *Boa) HelpFunc(cmd *cobra.Command, s []string) {
	model := newCmdModel(b.options, cmd)
	if err := tea.NewProgram(model, b.options.altScreen, b.options.mouseCellMotion).Start(); err != nil {
		log.Fatal(err)
	}
	if model.print {
		fmt.Println(b.options.styles.Border.Render(b.options.styles.CmdPrint.Render(model.cmdChain)))
	}
}

func (b *Boa) UsageFunc(cmd *cobra.Command) error {
	model := newCmdModel(b.options, cmd)
	if err := tea.NewProgram(model, b.options.altScreen, b.options.mouseCellMotion).Start(); err != nil {
		return err
	}
	if model.print {
		fmt.Println(b.options.styles.Border.Render(b.options.styles.CmdPrint.Render(model.cmdChain)))
	}
	return nil
}
