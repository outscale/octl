package filter

import (
	"context"
	"fmt"
	"iter"

	"github.com/charmbracelet/lipgloss"
	"github.com/itchyny/gojq"
	"github.com/outscale/octl/pkg/output/result"
)

const (
	baseColorIndex = 39
	colorOffset    = 16
)

type Color struct {
	query  *gojq.Query
	groups map[string]lipgloss.Style
}

func NewColor(s string) (*Color, error) {
	query, err := gojq.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid color expression: %w", err)
	}
	return &Color{
		query:  query,
		groups: make(map[string]lipgloss.Style),
	}, nil
}

func (c *Color) Filter(ctx context.Context, seq iter.Seq[result.Result]) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		for v := range seq {
			if v.Error != nil {
				_ = yield(v)
				return
			}
			iter := c.query.RunWithContext(ctx, v.Ok)
			jgroup, ok := iter.Next()
			if ok {
				if err, ok := jgroup.(error); ok {
					if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil { //nolint
						goto YIELD
					}
					_ = yield(result.Result{Error: fmt.Errorf("color error: %w", err)})
					return
				}
				v.Color = c.color(fmt.Sprint(jgroup))
			}
		YIELD:
			if !yield(v) {
				return
			}
		}
	}
}

func (c *Color) color(group string) lipgloss.Style {
	if st, ok := c.groups[group]; ok {
		return st
	}
	st := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(baseColorIndex + colorOffset*len(c.groups)))
	c.groups[group] = st
	return st
}

var _ Interface = (*Color)(nil)
