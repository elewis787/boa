package boa

import (
	tea "github.com/charmbracelet/bubbletea"
)

type options struct {
	// public
	altScreen tea.ProgramOption
	styles    *Styles
	// private (not capable of being set)
	mouseCellMotion tea.ProgramOption
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
			opt.altScreen = noOpt
			return
		}
		opt.altScreen = tea.WithAltScreen()
	})
}

func WithStyles(styles *Styles) Options {
	return newFuncOption(func(opt *options) {
		if styles != nil {
			opt.styles = styles
		}
	})
}

func defaultOptions() *options {
	return &options{
		altScreen:       noOpt,
		styles:          DefaultStyles(),
		mouseCellMotion: tea.WithMouseCellMotion(),
	}
}

func noOpt(*tea.Program) {}
