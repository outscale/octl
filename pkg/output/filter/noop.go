/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package filter

import (
	"context"
	"iter"

	"github.com/outscale/octl/pkg/output/result"
)

type Noop struct{}

func NewNoop() *Noop {
	return &Noop{}
}

func (e Noop) Filter(ctx context.Context, seq iter.Seq[result.Result]) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		for v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

var _ Interface = Noop{}
