/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package filter

import (
	"context"
	"fmt"
	"iter"

	"github.com/outscale/octl/pkg/output/result"
	"github.com/outscale/octl/pkg/structs"
)

type JSON struct{}

func (j JSON) Filter(ctx context.Context, seq iter.Seq[result.Result]) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		for v := range seq {
			if v.Error != nil {
				_ = yield(v)
				return
			}
			u, err := structs.ToUnstructured(v.Ok)
			if err != nil {
				_ = yield(result.Result{Error: fmt.Errorf("to unstructured: %w", err)})
				return
			}
			if !yield(result.New(v, u)) {
				return
			}
		}
	}
}

var _ Interface = JSON{}
