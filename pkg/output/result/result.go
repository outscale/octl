package result

import (
	"github.com/charmbracelet/lipgloss"
)

type Result struct {
	// Iter defines the iteration counter in a --watch/--wait call
	Iter        int
	SingleEntry bool
	// Profile defines the name of the profile the result comes from
	Profile string
	Ok      any
	Error   error
	Style   lipgloss.Style
}

func Error(err error) Result {
	return Result{Error: err}
}

func New(old Result, newOk any) Result {
	return Result{Ok: newOk, Iter: old.Iter, SingleEntry: old.SingleEntry, Profile: old.Profile}
}

func (r Result) WithStyle(s lipgloss.Style) Result {
	r.Style = s
	return r
}
