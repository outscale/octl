/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package filter

import (
	"context"
	"iter"
	"time"

	"github.com/outscale/octl/pkg/output/result"
)

type Elapsed struct {
	start time.Time
}

func NewElapsed() *Elapsed {
	return &Elapsed{
		start: time.Now(),
	}
}

func (e Elapsed) Filter(ctx context.Context, seq iter.Seq[result.Result]) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		for v := range seq {
			if v.Error != nil {
				_ = yield(v)
				return
			}
			if json, ok := v.Ok.(map[string]any); ok {
				dur := time.Since(e.start)
				if dur > time.Second {
					dur = dur.Truncate(time.Second)
				} else {
					dur = dur.Truncate(time.Millisecond)
				}
				json["_elapsed"] = dur.String()
			}
			if !yield(v) {
				return
			}
		}
	}
}

var _ Interface = Elapsed{}
