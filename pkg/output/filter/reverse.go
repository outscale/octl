/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package filter

import (
	"context"
	"iter"
	"slices"

	"github.com/outscale/octl/pkg/output/result"
)

type Reverse struct{}

func (j Reverse) Filter(ctx context.Context, seq iter.Seq[result.Result]) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		res := slices.Collect(seq)
		slices.Reverse(res)
		for _, v := range res {
			if !yield(v) {
				return
			}
		}
	}
}

var _ Interface = Reverse{}
