/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package filter

import (
	"context"
	"iter"
	"strings"

	"github.com/outscale/octl/pkg/config"
	"github.com/outscale/octl/pkg/output/format"
	"github.com/outscale/octl/pkg/output/result"
	"github.com/samber/lo"
)

type Disappear struct {
	cols, colState int
	iter           int
	last           map[string]struct{}
	current        map[string]struct{}
}

func (d *Disappear) Filter(ctx context.Context, seq iter.Seq[result.Result]) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		for v := range seq {
			if v.Error != nil {
				_ = yield(v)
				return
			}
			if v.Iter != d.iter {
				// checking which entries were present in the last batch but not the current
				for id := range d.last {
					if _, found := d.current[id]; !found {
						row := make([]string, d.cols)
						row[0] = id
						row[d.colState] = format.StateDisappeared
						if !yield(result.New(v, row)) {
							return
						}
					}
				}
				d.iter = v.Iter
				d.current, d.last = d.last, d.current
				clear(d.current)
			}

			if row, ok := v.Ok.([]string); ok && len(row) > 0 {
				d.current[row[0]] = struct{}{}
			}
			if !yield(v) {
				return
			}
		}
	}
}

func NewDisappear(cols config.Columns) Interface {
	if len(cols) == 0 {
		return NewNoop()
	}
	suffix := lo.LastOrEmpty(lo.Words(cols[0].Title))
	if !strings.EqualFold(suffix, "id") && !strings.EqualFold(suffix, "name") {
		return NewNoop()
	}
	_, colState, found := lo.FindIndexOf(cols, func(c config.Column) bool {
		suffix := lo.LastOrEmpty(lo.Words(c.Title))
		return strings.EqualFold(suffix, "state") || strings.EqualFold(suffix, "status")
	})
	if !found {
		return NewNoop()
	}
	return &Disappear{
		cols:     len(cols),
		colState: colState,
		last:     make(map[string]struct{}),
		current:  make(map[string]struct{}),
	}
}

var _ Interface = (*Dedup)(nil)
