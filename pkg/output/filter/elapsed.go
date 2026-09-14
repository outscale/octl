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
		var dur time.Duration
		iter := -1
		for v := range seq {
			if v.Error != nil {
				_ = yield(v)
				return
			}
			if v.Iter != iter {
				dur = time.Since(e.start)
				if dur > time.Second {
					dur = dur.Truncate(time.Second)
				} else {
					dur = dur.Truncate(time.Millisecond)
				}
			}
			switch val := v.Ok.(type) {
			case map[string]any:
				val["_elapsed"] = dur.String()
			case string:
				if v.Iter != iter {
					if !yield(result.New(v, "--- elapsed since start: "+dur.String()+" ---")) {
						return
					}
				}
			}
			if !yield(v) {
				return
			}
			iter = v.Iter
		}
	}
}

var _ Interface = Elapsed{}
