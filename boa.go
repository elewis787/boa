package boa

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type Boa struct {
	options *options
	Styles  *Styles
}

// Create a new instance of Boa with custom Options
func New(options ...Options) *Boa {
	opts := defaultOptions()
	for _, opt := range options {
		opt.apply(opts)
	}
	return &Boa{
		Styles:  defaultStyles(opts.width),
		options: opts,
	}
}

func (b *Boa) HelpFunc(cmd *cobra.Command, s []string) {
	model := newCmdModel(b.options, b.Styles, cmd)
	if err := tea.NewProgram(model, b.options.atlScreen, b.options.mouseCellMotion).Start(); err != nil {
		log.Fatal(err)
	}
	if model.print {
		fmt.Println(model.cmdChain)
	}
}

func (b *Boa) UsageFunc(cmd *cobra.Command) error {
	model := newCmdModel(b.options, b.Styles, cmd)
	if err := tea.NewProgram(model, b.options.atlScreen, b.options.mouseCellMotion).Start(); err != nil {
		return err
	}
	if model.print {
		fmt.Println(model.cmdChain)
	}
	return nil
}
