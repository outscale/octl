/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package watch

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/outscale/octl/pkg/output/format"
)

type Format struct {
	mu    sync.Mutex
	ctx   context.Context
	v     any
	fmter format.Interface
}

func NewFormat(fmter format.Interface) *Format {
	return &Format{
		fmter: fmter,
	}
}

// Init implements tea.Model.
func (m *Format) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m *Format) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// View implements tea.Model.
func (m *Format) View() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.v == nil {
		return ""
	}
	var buf bytes.Buffer
	_ = m.fmter.Format(m.ctx, &buf, m.v)
	return buf.String()
}

// Format implements format.Interface.
// It does not display anything, as output is managed by Bubbletea through the View() method.
func (m *Format) Format(ctx context.Context, w io.Writer, v any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ctx = ctx
	vv := reflect.ValueOf(v)
	if m.v == nil {
		m.v = reflect.MakeSlice(vv.Type(), 0, 5).Interface()
	}
	// if the input was a single entry, make it a slice
	if vv.Kind() != reflect.Slice {
		sv := reflect.MakeSlice(vv.Type(), 1, 1)
		sv.Index(0).Set(vv)
		vv = sv
	}
	mv := reflect.ValueOf(m.v)
	for i := range vv.Len() {
		mv = reflect.Append(mv, vv.Index(i))
	}
	m.v = mv.Interface()
	return nil
}

// Error implements format.Interface.
func (m *Format) Error(ctx context.Context, v any) error {
	return m.fmter.Error(ctx, v)
}

var (
	_ tea.Model        = (*Format)(nil)
	_ format.Interface = (*Format)(nil)
)
