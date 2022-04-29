package boa

import tea "github.com/charmbracelet/bubbletea"

type options struct {
	atlScreen tea.ProgramOption
}

type Options interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

func (fo *funcOption) apply(opt *options) {
	fo.f(opt)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{f: f}
}

func WithAltScreen(b bool) Options {
	return newFuncOption(func(opt *options) {
		if !b {
			opt.atlScreen = noOpt
		}
		opt.atlScreen = tea.WithAltScreen()
	})
}

func defaultOptions() *options {
	return &options{
		atlScreen: noOpt,
	}
}

func noOpt(*tea.Program) {}

// Inital options
var opts = defaultOptions()

// New uses sets the package level options variables.
// Eventually this will likely return a new struct and
// the options will belong to that struct
func New(options ...Options) {
	for _, opt := range options {
		opt.apply(opts)
	}
}
