package boa

import tea "github.com/charmbracelet/bubbletea"

type options struct {
	// public
	altScreen  tea.ProgramOption
	width      int
	showBorder bool

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

func WithBorder(b bool) Options {
	return newFuncOption(func(opt *options) {
		opt.showBorder = b
	})
}

func WithWidth(w int) Options {
	return newFuncOption(func(opt *options) {
		opt.width = w
	})
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

func defaultOptions() *options {
	return &options{
		width:           defaultWidth,
		altScreen:       noOpt,
		showBorder:      true,
		mouseCellMotion: tea.WithMouseCellMotion(),
	}
}

func noOpt(*tea.Program) {}
