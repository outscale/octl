package result

import "github.com/charmbracelet/lipgloss"

type Result struct {
	SingleEntry bool
	Ok          any
	Error       error
	Color       lipgloss.Style
}
